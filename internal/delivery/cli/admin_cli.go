package cli

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/SergeyBogomolovv/fitflow/internal/domain"
)

type AdminService interface {
	CreateAdmin(ctx context.Context, login, password string) error
	UpdatePassword(ctx context.Context, login, oldPass, newPass string) error
	RemoveAdmin(ctx context.Context, login string) error
}

type adminCli struct {
	reader *bufio.Reader
	svc    AdminService
}

func NewAdminCLI(svc AdminService) *adminCli {
	reader := bufio.NewReader(os.Stdin)
	return &adminCli{reader, svc}
}

func (c *adminCli) Run(ctx context.Context) {
	if len(os.Args) < 2 {
		fmt.Println("Ожидается команда: create, update-password, remove")
		return
	}
	command := os.Args[1]

	switch command {
	case "create":
		c.handleCreate(ctx)
	case "update-password":
		c.handleUpdatePassword(ctx)
	case "remove":
		c.handleRemove(ctx)
	default:
		fmt.Println("Неизвестная команда. Используйте: create, update-password, remove")
	}
}

func (c *adminCli) handleCreate(ctx context.Context) {
	login := c.readPrompt("Логин: ")
	password := c.readPrompt("Пароль: ")
	if login == "" || password == "" {
		fmt.Println("Логин и пароль не могут быть пустыми.")
		return
	}
	err := c.svc.CreateAdmin(ctx, login, password)
	if err != nil {
		if errors.Is(err, domain.ErrAdminAlreadyExists) {
			fmt.Println("Администратор с таким логином уже существует.")
			return
		}
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Println("Администратор создан.")
}

func (c *adminCli) handleUpdatePassword(ctx context.Context) {
	login := c.readPrompt("Логин: ")
	oldPass := c.readPrompt("Старый пароль: ")
	password := c.readPrompt("Новый пароль: ")
	err := c.svc.UpdatePassword(ctx, login, oldPass, password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			fmt.Println("Неверные данные")
			return
		}
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Println("Пароль обновлён.")
}

func (c *adminCli) handleRemove(ctx context.Context) {
	login := c.readPrompt("Логин: ")
	err := c.svc.RemoveAdmin(ctx, login)
	if err != nil {
		if errors.Is(err, domain.ErrAdminNotFound) {
			fmt.Println("Администратор не существует.")
		}
		fmt.Printf("Ошибка: %v\n", err)
		return
	}
	fmt.Println("Администратор удалён.")
}

func (c *adminCli) readPrompt(prompt string) string {
	fmt.Print(prompt)
	input, err := c.reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Ошибка ввода: %v", err)
	}
	return strings.TrimSpace(input)
}
