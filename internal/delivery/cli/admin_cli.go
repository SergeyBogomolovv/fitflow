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
	svc AdminService
}

func NewAdminCLI(svc AdminService) *adminCli {
	return &adminCli{svc}
}

func (c *adminCli) Run(ctx context.Context) {
	if len(os.Args) < 2 {
		fmt.Println("Ожидается команда: create, update-password, remove")
		return
	}
	command := os.Args[1]

	reader := bufio.NewReader(os.Stdin)

	switch command {
	case "create":
		c.handleCreate(ctx, reader)
	case "update-password":
		c.handleUpdatePassword(ctx, reader)
	case "remove":
		c.handleRemove(ctx, reader)
	default:
		fmt.Println("Неизвестная команда. Используйте: create, update-password, remove")
	}
}

func (c *adminCli) handleCreate(ctx context.Context, reader *bufio.Reader) {
	login := readPrompt(reader, "Логин: ")
	password := readPrompt(reader, "Пароль: ")
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

func (c *adminCli) handleUpdatePassword(ctx context.Context, reader *bufio.Reader) {
	login := readPrompt(reader, "Логин: ")
	oldPass := readPrompt(reader, "Старый пароль: ")
	password := readPrompt(reader, "Новый пароль: ")
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

func (c *adminCli) handleRemove(ctx context.Context, reader *bufio.Reader) {
	login := readPrompt(reader, "Логин: ")
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

func readPrompt(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Ошибка ввода: %v", err)
	}
	return strings.TrimSpace(input)
}
