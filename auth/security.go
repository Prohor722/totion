package auth

import "github.com/Prohor722/totion/util"

func (s *authService) isAccountLocked(username string) bool {
    s.mutex.RLock()
    defer s.mutex.RUnlock()
    failure, exists := s.failures[username]
    if !exists {
        return false
    }
    return s.clock.Now().Before(failure.LockedUntil)
}

func (s *authService) recordFailedLogin(username string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    failure := s.failures[username]
    now := s.clock.Now()
    if now.Sub(failure.LastAttempt) > util.FailedLoginWindow {
        failure.Count = 0
    }
    failure.Count++
    failure.LastAttempt = now
    if failure.Count >= util.FailedLoginThreshold {
        failure.LockedUntil = now.Add(util.AccountLockDuration)
    }
    s.failures[username] = failure
}

func (s *authService) clearFailedLogin(username string) {
    s.mutex.Lock()
    defer s.mutex.Unlock()
    delete(s.failures, username)
}
