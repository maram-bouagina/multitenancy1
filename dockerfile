# Télécharger Go pour la compilation
FROM golang:1.24-alpine

# Dossier de travail
WORKDIR /app

# Copier la liste des dépendances
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier le code source
COPY . .

# Compiler le code en un fichier exécutable appelé api
RUN go build -o api ./main.go

# Port exposé
EXPOSE 8000

# Lancer le programme
CMD ["./api"]