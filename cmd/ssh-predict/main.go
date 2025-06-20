package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/GOODMANEXPERT/ssh-predict/checks"
	"github.com/GOODMANEXPERT/ssh-predict/internal/platform"
	"github.com/GOODMANEXPERT/ssh-predict/parser"
	"github.com/GOODMANEXPERT/ssh-predict/report"
)

func main() {
	// 1) Чтение флага --config
	cfgPath := flag.String("config", "", "путь к sshd_config")
	flag.Parse()

	// 2) Если флаг пуст, получить путь по ОС
	if *cfgPath == "" {
		p, err := platform.ConfigPath()
		if err != nil {
			fmt.Println("[ERROR] ОС не поддерживается:", err)
			os.Exit(1)
		}
		cfgPath = &p
	}

	// 3) Парсинг конфигурации
	cp := parser.NewConfigParser(*cfgPath)
	cfg, err := cp.Parse()
	if err != nil {
		fmt.Printf("[ERROR] Не удалось прочитать %s: %v\n", *cfgPath, err)
		os.Exit(1)
	}

	// 4) Список проверок
	checkers := []checks.Checker{
		checks.PermitRootLoginChecker(),
	}

	// 5) Выполнение проверок
	var results []checks.CheckResult
	for _, c := range checkers {
		results = append(results, c.Check(cfg))
	}

	// 6) Вывод отчёта
	reporter := report.NewStdoutReporter()
	reporter.Report(results)
}
