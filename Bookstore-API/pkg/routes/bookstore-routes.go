package routes

import (
  "github.com/gorilla/mux"
  "github.com/Ridwan-Al-Mahmud/Go-Bookstore/pkg/controllers"
)

var RegisterBookStoreRoutes = func(router *mux.Router){
  router.HandleFunc("/book/", controllers.CreateBook).Methods("POST")
  router.HandleFunc("/book/", controllers.GetBooks).Methods("GET")
  router.HandleFunc("/book/{bookId}", controllers.GetBooksById).Methods("GET")
  router.HandleFunc("/book/{bookId}", controllers.UpdateBook).Methods("PUT")
  router.HandleFunc("/book/{bookId}", controllers.DeleteBook).Methods("DELETE")
}