package auth

import(
	"context"
	"errors"
	"testing"
)



func TestRegister_WeakPass(t *testing.T){
	service := NewService(nil)


	_, err:= service.Register(context.Background(), "ana@gmail.com", "abc", "Ana")


	if !errors.Is(err, ErrWeakPassword) {
	t.Errorf("got %v, want ErrWeakPassword", err)
}
}


 

