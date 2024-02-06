package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-crud/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)


type responce struct{
	ID 		int64	`json:"id,omitempty"`
	Message string `json:"message.omitempty"`
}

func createConnection() *sql.DB{
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading the .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("SUccessfully connected to postgres database")
	return db
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertId := insertUser(user)

	res := responce{
		ID: insertId,
		Message: "User created successfully",
	}

	json.NewEncoder(w).Encode(res)

}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi((params["id"]))

	if err != nil {
		log.Fatalf("Unable to conert the string onto integer. %v", err)
	}

	user, err := getUser(int64(id))
	if err != nil {
		log.Fatalf("unable to get %v", err)
	}
	json.NewEncoder(w).Encode(user)
}

func GetAllUser(w http.ResponseWriter, r *http.Request){
	user, err := getAllUsers()

	if err != nil {
		log.Fatalf("Unable to getall the users%v", err)
	}

	json.NewEncoder(w).Encode(user)

	
}

func UpdateUser(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	var user models.User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatalf("Unable to decode the requires body. %v", err)
	}

	updatedRows := updateUser(int64(id), user)

	msg := fmt.Sprintf("user updated. totola rows affected %v", updatedRows)

	res := responce {
		ID: int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func DeleteUser(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert to int %v", err)
	}

	deleteRow := deleteUser(int64(id))
	msg := fmt.Sprintf("User deleted successfully, row affected %v", deleteRow)

	res := responce{
		ID: int64(id),
		Message: msg,
	}

	json.NewEncoder(w).Encode(res)
}


func insertUser(user models.User) int64{
	db := createConnection()

	defer db.Close()
	sqlStatement :=`INSERT INTO crud_users (name, age, email) VALUES ($1, $2, $3) RETURNING userid`
	var id int64

	err := db.QueryRow(sqlStatement, user.Name, user.Age, user.Email).Scan(&id)
	if err != nil {
		log.Fatalf("Unabel to execute the query, %v", err)
	}

	fmt.Printf("incerted a single record %d", id)
	return id
}

func getUser(id int64)(models.User, error){
	db := createConnection()

	defer db.Close()

	var user models.User

	sqlStatement := `SELECT * FROM crud_users WHERE userid=$1`

	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&user.UserId, &user.Name, &user.Age, &user.Email)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Now row were returned")
		return user, nil
	case nil :
		return user, nil
	default:
		log.Fatal("unable to get the row")
	}

	return user, err

}

func getAllUsers()([]models.User, error) {
	db := createConnection()

	defer db.Close()

	var users []models.User

	sqlStatement := `SELECT * FROM crud_users`

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("inable to execite the query %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var user models.User
		err = rows.Scan(&user.UserId, &user.Name, &user.Age, &user.Email)

		if err != nil {
			log.Fatalf("Unable to scan the row %v", err)
		}
		users = append(users, user)
	}

	return users, err

}

func updateUser(id int64, user models.User) int64{
	db := createConnection()

	defer db.Close()

	sqlStatement := `Update crud_users SET name=$2, age=$3, email=$4 WHERE userid=$1`

	res, err := db.Exec(sqlStatement, id, user.Name, user.Age, user.Email)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking he affected row, %v", err)
	}

	fmt.Printf("Total sffected rows %v", rowsAffected)

	return rowsAffected


}

func deleteUser(id int64) int64 {
	db := createConnection()

	defer db.Close()

	sqlStatement := `DELETE FROM crud_users WHERE userid=$1`
	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking he affected row, %v", err)
	}

	fmt.Printf("Total sffected rows %v", rowsAffected)

	return rowsAffected


}	