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

func TestRegister_InvalidEmail(t *testing.T){
	service := NewService(nil)


	_, err:= service.Register(context.Background(), "ognen", "abcbc2222", "Ana")


	if !errors.Is(err, ErrInvalidEmail) {
	t.Errorf("got %v, want ErrInvalidEmail", err)
}
}


func TestRegister_LongPassword(t *testing.T){
	service := NewService(nil)

	_, err := service.Register(context.Background(), "ana@gmail.com", "abcaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa" , "Ana")


if !errors.Is(err, ErrWeakPassword ){
	t.Errorf("got %v, want ErrWeakPassword", err)
}
}
func TestRegister_EmptyName(t *testing.T){
	service := NewService(nil)


	_, err:= service.Register(context.Background(), "ana@gmail.com", "abcbc222", "    ")


	if !errors.Is(err, ErrNameRequired) {
	t.Errorf("got %v, want ErrNameRequired", err)
}
}
 

