package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	db "github.com/mktsy/go-login-api/api/databases"
	middlewares "github.com/mktsy/go-login-api/api/middlewares"
	model "github.com/mktsy/go-login-api/api/models"
)

//GetUser get a user by Id
func GetUser(s *db.Dispatch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ss := s.MongoDB.Copy()
		defer ss.Close()

		// Get id from param
		id := chi.URLParam(r, "id")

		// Verify id is ObjectId, otherwise bail
		if !bson.IsObjectIdHex(id) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		oid := bson.ObjectIdHex(id)
		u := model.User{}
		if err := ss.DB("login").C("users").FindId(oid).One(&u); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		uj, _ := json.Marshal(u)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", uj)
	}
}

//CreateUser create a new user
func CreateUser(s *db.Dispatch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ss := s.MongoDB.Copy()
		defer ss.Close()

		// Stub an user to be populated from the body
		u := model.User{}
		json.NewDecoder(r.Body).Decode(&u)

		// Add an Id
		u.Id = bson.NewObjectId()
		u.CreatedAt = time.Now()
		u.UpdatedAt = time.Now()

		if passwd, err := middlewares.Encrypt(u.Password); err == nil {
			u.Password = passwd
		}

		ss.DB("login").C("users").Insert(u)
		uj, _ := json.Marshal(u)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", uj)
	}
}

// DeleteUser remove user from database
func DeleteUser(s *db.Dispatch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ss := s.MongoDB.Copy()
		defer ss.Close()

		id := chi.URLParam(r, "id")

		if !bson.IsObjectIdHex(id) {
			msg := []byte(`{"message":"ObjectId invalid"}`)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s", msg)
			return
		}

		c := ss.DB("login").C("users")

		if err := c.Remove(bson.M{"_id": bson.ObjectIdHex(id)}); err != nil {
			switch err {
			default:
				msg := []byte(`{"message":"ObjectId invalid"}`)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", msg)

			case mgo.ErrNotFound:
				msg := []byte(`{"message":"ObjectId not found"}`)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", msg)
			}
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateUser(s *db.Dispatch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ss := s.MongoDB.Copy()
		defer ss.Close()

		id := chi.URLParam(r, "id")

		// Verify id is ObjectId, otherwise bail
		if !bson.IsObjectIdHex(id) {
			msg := []byte(`{"message":"ObjectId invalid"}`)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s", msg)
			return
		}

		// Stub an user to be populated from the body
		u := model.User{}
		json.NewDecoder(r.Body).Decode(&u)
		u.UpdatedAt = time.Now()
		u.Id = bson.ObjectIdHex(id)

		// if err := ss.DB("login").C("users").FindId(u.Id).One(&u); err != nil {
		// 	w.WriteHeader(http.StatusNotFound)
		// 	return
		// }

		// uj, _ := json.Marshal(u)
		// fmt.Printf("%s", uj)

		c := ss.DB("login").C("users")

		if err := c.Update(bson.M{"_id": u.Id}, &u); err != nil {
			switch err {
			default:
				msg := []byte(`{"message":"ObjectId invalid"}`)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Print(id)
				fmt.Fprintf(w, "%s", msg)
			case mgo.ErrNotFound:
				msg := []byte(`{"message":"ObjectId not found"}`)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "%s", msg)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
