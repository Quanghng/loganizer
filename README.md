# 🛠️ GoLog Analyzer (`loganizer`)

Un outil en ligne de commande (CLI) écrit en Go qui permet d’analyser des fichiers de logs de manière **concurrente**, avec gestion d’erreurs personnalisées et export des résultats en **JSON**.

Ce projet est développé dans le cadre du TP « Analyse de Logs Distribuée » et met en pratique les concepts de **goroutines**, **WaitGroups**, **canaux (channels)**, **erreurs personnalisées**, et **Cobra** pour la création d’outils CLI.

---

## 🚀 Fonctionnalités

* **Analyse de logs concurrente :**

  * Chaque fichier est analysé dans une goroutine indépendante.
  * Simulation de l’analyse avec un délai aléatoire (`50–200ms`).

* **Gestion robuste des erreurs :**

  * `InaccessibleFileError` → fichier introuvable ou non lisible.
  * `ParsingError` → problème de format du log.

* **Export JSON :**

  * Résultats collectés dans un tableau JSON.
  * Export optionnel via le flag `--output` (`-o`).
  * Champs exportés : `log_id`, `file_path`, `status`, `message`, `error_details`.

* **CLI basée sur Cobra :**

  * Commande principale : `analyze`.
  * Flags configurables : `--config` et `--output`.

* **Affichage console en temps réel :**

  * Chaque log est affiché avec son statut (`✅ OK`, `❌ FAILED`, etc.).

---

## 📂 Structure du projet

```
.
├── cmd/
│   ├── root.go         # Commande racine Cobra
│   └── analyze.go      # Implémentation de la commande `analyze`
├── internal/
│   ├── analyzer/       # Analyse, erreurs personnalisées, rapport
│   ├── config/         # Lecture du fichier config.json
│   └── reporter/       # Export JSON des résultats
├── config.json         # Exemple de fichier de configuration
├── go.mod
├── go.sum
└── main.go             # Point d’entrée du programme
```

---

## ⚙️ Prérequis

* [Go 1.23+](https://go.dev/dl/)
* Cobra installé automatiquement via `go get github.com/spf13/cobra@latest`
* Windows PowerShell (les commandes ci-dessous utilisent `.\bin\loganizer.exe`)

---

## 📦 Installation

Clonez le dépôt puis construisez l’application :

```powershell
git clone https://github.com/Quanghng/loganizer
cd loganizer
go build -o ./bin/loganizer.exe .
```

---

## ▶️ Utilisation

### 1. Préparer un fichier de configuration

Créez un fichier `config.json` contenant les logs à analyser :

```json
[
  {
    "id": "web-server-1",
    "path": "C:/logs/nginx/access.log",
    "type": "nginx-access"
  },
  {
    "id": "app-backend-2",
    "path": "C:/logs/app/errors.log",
    "type": "custom-app"
  },
  {
    "id": "invalid-path",
    "path": "C:/invalid/missing.log",
    "type": "unknown"
  }
]
```

### 2. Lancer l’analyse

Commande de base :

```powershell
.\bin\loganizer.exe analyze -c "./config.json"
```

Avec export JSON :

```powershell
.\bin\loganizer.exe analyze -c "./config.json" -o "./report.json"
```

### 3. Exemple de sortie console

```
✅ web-server-1 (path: test_logs/access.log) : OK
✅ app-backend-2 (path: test_logs/errors.log) : OK
✅ corrupted-log (path: test_logs/corrupted.log) : OK
🚫 db-server-3 est inaccessible : GetFileAttributesEx test_logs/mysql_error.log: The system cannot find the file specified.
🚫 invalid-path est inaccessible : GetFileAttributesEx /non/existent/log.log: The system cannot find the path specified.
✅ empty-log (path: test_logs/empty.log) : OK
✅ Résultats exportés vers ./report.json
```

### 4. Exemple de `report.json`

```json
[
  {
    "log_id": "web-server-1",
    "file_path": "C:/logs/nginx/access.log",
    "status": "OK",
    "message": "Analyse terminée avec succès.",
    "error_details": ""
  },
  {
    "log_id": "invalid-path",
    "file_path": "C:/invalid/missing.log",
    "status": "FAILED",
    "message": "Fichier introuvable.",
    "error_details": "file inaccessible: open C:/invalid/missing.log: The system cannot find the path specified."
  }
]
```

---

## 📖 Documentation interne

* `internal/config` : contient `LoadTargetsFromFile`, fonction qui lit un JSON et retourne une liste de `InputTarget`.
* `internal/analyzer` :

  * `AnalyzeLogSync` : simule l’analyse d’un fichier et retourne un `AnalyzeResult`.
  * `ConvertToReportEntry` : convertit un `AnalyzeResult` en `ReportEntry`.
  * `InaccessibleFileError`, `ParsingError` : erreurs personnalisées.
* `internal/reporter` :

  * `ExportResultsToJsonFile` : sérialise une slice de `ReportEntry` dans un fichier JSON.
* `cmd/` : commandes Cobra (`root`, `analyze`).

---

## Améliorations

* Création automatique des dossiers de sortie avec `os.MkdirAll`.
* Ajout d’un horodatage dans le nom du rapport (`240924_report.json`).
* Nouvelle sous-commande `add-log` pour ajouter une config au fichier JSON.
* Filtrage des résultats avec `--status OK` ou `--status FAILED`.

