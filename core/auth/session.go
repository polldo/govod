package auth

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/polldo/govod/api/web"
	"github.com/polldo/govod/api/weberr"
	"github.com/polldo/govod/core/claims"
)

const userKey = "userID"
const roleKey = "role"

func Authenticate(s *scs.SessionManager) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			uid, ok := s.Get(ctx, userKey).(string)
			if !ok {
				return weberr.NotAuthorized(errors.New("no userID in session"))
			}

			role, ok := s.Get(ctx, roleKey).(string)
			if !ok {
				return weberr.NotAuthorized(errors.New("no user role in session"))
			}

			ctx = claims.Set(ctx, claims.Claims{UserID: uid, Role: role})

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}

func Admin(s *scs.SessionManager) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			role, ok := s.Get(ctx, roleKey).(string)
			if !ok {
				return weberr.NotAuthorized(errors.New("no user role in session"))
			}

			if role != claims.RoleAdmin {
				return weberr.NotAuthorized(fmt.Errorf("user role is not admin: %s", role))
			}

			return handler(ctx, w, r)
		}
		return h
	}
	return m
}

func LoadAndSave(s *scs.SessionManager) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			var token string
			cookie, err := r.Cookie(s.Cookie.Name)
			if err == nil {
				token = cookie.Value
			}

			ctx, err = s.Load(ctx, token)
			if err != nil {
				return err
			}

			bw := &bufferedResponseWriter{ResponseWriter: w}
			if err := handler(ctx, bw, r); err != nil {
				return err
			}

			if r.MultipartForm != nil {
				r.MultipartForm.RemoveAll()
			}

			switch s.Status(ctx) {
			case scs.Modified:
				token, expiry, err := s.Commit(ctx)
				if err != nil {
					return err
				}

				s.WriteSessionCookie(ctx, w, token, expiry)
			case scs.Destroyed:
				s.WriteSessionCookie(ctx, w, "", time.Time{})
			}

			w.Header().Add("Vary", "Cookie")

			if bw.code != 0 {
				w.WriteHeader(bw.code)
			}
			w.Write(bw.buf.Bytes())

			return nil
		}
		return h
	}
	return m
}

type bufferedResponseWriter struct {
	http.ResponseWriter
	buf         bytes.Buffer
	code        int
	wroteHeader bool
}

func (bw *bufferedResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func (bw *bufferedResponseWriter) WriteHeader(code int) {
	if !bw.wroteHeader {
		bw.code = code
		bw.wroteHeader = true
	}
}

func (bw *bufferedResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj := bw.ResponseWriter.(http.Hijacker)
	return hj.Hijack()
}

func (bw *bufferedResponseWriter) Push(target string, opts *http.PushOptions) error {
	if pusher, ok := bw.ResponseWriter.(http.Pusher); ok {
		return pusher.Push(target, opts)
	}
	return http.ErrNotSupported
}
