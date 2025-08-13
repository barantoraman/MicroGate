package transport

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/barantoraman/microgate/internal/auth/pb"
	"github.com/barantoraman/microgate/internal/auth/repo/entity"
	"github.com/barantoraman/microgate/pkg/config"
	contx "github.com/barantoraman/microgate/pkg/ctx"
	errs "github.com/barantoraman/microgate/pkg/err"
	"github.com/barantoraman/microgate/pkg/token"
	"github.com/barantoraman/microgate/pkg/validator"
	httpTransport "github.com/go-kit/kit/transport/http"
	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}
	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)
		mu.Lock()
		if _, found := clients[ip]; !found {
			clients[ip] = &client{
				limiter: rate.NewLimiter(
					rate.Limit(2),
					4,
				),
			}
		}
		clients[ip].lastSeen = time.Now()
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			errs.RateLimitExceededResponse(w)
			return
		}
		mu.Unlock()
		next.ServeHTTP(w, r)
	})
}

func authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = contx.SetUser(r, entity.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			w.Write([]byte(fmt.Sprintf("%s\n", headerParts)))
			errs.InvalidAuthenticationToken(w)
			return
		}

		tkn := headerParts[1]

		v := validator.New()

		if token.ValidateTokenPlaintext(v, tkn); !v.Valid() {
			errs.InvalidAuthenticationToken(w)
			return
		}

		var cfg config.ApiGatewayServiceConfigurations
		err := config.GetLoader().GetConfigByKey("api_gateway_service", &cfg)
		if err != nil {
		}

		var grpcAddr = net.JoinHostPort(cfg.AuthServiceHost, cfg.AuthServicePort)
		conn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("fail to dial: %v", err)
		}
		defer conn.Close()

		client := pb.NewAuthClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		pToken := &pb.Token{
			PlaintText: tkn,
			Hash:       nil,
			UserId:     0,
			Expiry:     nil,
			Scope:      token.ScopeAuthentication,
		}

		resp, err := client.IsAuth(ctx, &pb.IsAuthRequest{Token: pToken})
		if err != nil {
			errs.InvalidAuthenticationToken(w)
			return
		}

		usr := &entity.User{
			UserID:       resp.Token.UserId,
			Email:        "",
			Password:     "",
			PasswordHash: nil,
		}
		r = contx.SetUser(r, usr)

		next.ServeHTTP(w, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Frame-Options", "deny")

		next.ServeHTTP(w, r)
	})
}

func requireAuthenticatedUser(next *httpTransport.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := contx.GetUser(r)

		if user.IsAnonymous() {
			errs.AuthenticationRequiredResponse(w)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Origin")

		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := w.Header().Get("Origin")

		if origin == "localhost:8080" {
			w.Header().Set("Access-Control-Allow-Origin", origin)

			if r.Method == http.MethodOptions &&
				w.Header().Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods",
					"OPTIONS, PUT, PATCH, DELETE")
				w.Header().Set("Access-Control-Allow-Headers",
					"Authorization, Content-Type")

				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				errs.ServerErroResponse(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
