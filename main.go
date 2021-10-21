package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	database "github.com/OutBoulderCounty/icc-database"
	forms "github.com/OutBoulderCounty/icc-forms"
	_ "github.com/go-sql-driver/mysql"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoOption struct {
	ID    primitive.ObjectID `bson:"_id"`
	Name  string             `bson:"name"`
	Index int                `bson:"index"`
}

type mongoElement struct {
	ID       primitive.ObjectID `bson:"_id"`
	Label    string             `bson:"label"`
	Type     string             `bson:"type"`
	Index    int                `bson:"index"`
	Required bool               `bson:"required"`
	Options  []mongoOption      `bson:"options"`
	Priority int                `bson:"priority"`
	Search   bool               `bson:"search"`
}

type mongoForm struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `bson:"name"`
	Elements []mongoElement     `bson:"elements"`
	Required bool               `bson:"required"`
	Live     bool               `bson:"live"`
}

func main() {
	truncate := flag.Bool("truncate", false, "determines whether to truncate existing data on the target")
	flag.Parse()

	// get forms data from prod mongodb deployment
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
	if err != nil {
		log.Fatal("Error connecting to MongoDB: " + err.Error())
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB: " + err.Error())
	}
	fmt.Println("MongoDB connection successful")
	defer client.Disconnect(ctx)
	mongoDB := client.Database("healthdir")
	coll := mongoDB.Collection("forms")
	filter := bson.M{}
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		log.Fatal("Failed to retrieve forms: " + err.Error())
	}
	var mongoForms []*mongoForm
	err = cursor.All(ctx, &mongoForms)
	if err != nil {
		log.Fatal("Failed to parse forms: " + err.Error())
	}

	db, err := database.Connect("dev")
	if err != nil {
		log.Fatal("Failed to connect to SQL database: " + err.Error())
	}
	if *truncate {
		truncateForms := "DELETE FROM forms"
		_, err := db.Exec(truncateForms)
		if err != nil {
			log.Fatal("Failed to truncate forms: " + err.Error())
		}
	}

	count := 0
	for i := 0; i < len(mongoForms); i++ {
		// manipulate forms data to be SQL compatible
		mongoForm := mongoForms[i]
		sqlForm := forms.Form{
			Name:     mongoForm.Name,
			Required: mongoForm.Required,
			Live:     mongoForm.Live,
		}
		// insert form and get ID back
		sqlStmt := fmt.Sprintf("INSERT INTO forms (name, required, live) VALUES (\"%s\", %v, %v)", sqlForm.Name, sqlForm.Required, sqlForm.Live)
		result, err := db.Exec(sqlStmt)
		if err != nil {
			fmt.Println("Failed SQL: " + sqlStmt)
			log.Fatal("Failed to insert form: " + err.Error())
		}
		id, err := result.LastInsertId()
		if err != nil {
			log.Fatal("Failed to get inserted form ID: " + err.Error())
		}

		// insert elements
		for j := 0; j < len(mongoForm.Elements); j++ {
			element := mongoForm.Elements[j]
			sqlElement := forms.Element{
				FormID:   id,
				Label:    element.Label,
				Type:     element.Type,
				Position: element.Index,
				Required: element.Required,
				Priority: element.Priority,
				Search:   element.Search,
			}
			insertElement := fmt.Sprintf("INSERT INTO elements (formID, label, type, position, required, priority, search) VALUES (%v, \"%s\", \"%s\", %v, %v, %v, %v)", sqlElement.FormID, sqlElement.Label, sqlElement.Type, sqlElement.Position, sqlElement.Required, sqlElement.Priority, sqlElement.Search)
			elemResult, err := db.Exec(insertElement)
			if err != nil {
				fmt.Println("Failed SQL: " + insertElement)
				log.Fatal("Failed to insert element: " + err.Error())
			}
			elemId, err := elemResult.LastInsertId()
			if err != nil {
				log.Fatal("Failed to get inserted element ID: " + err.Error())
			}
			// insert options
			for k := 0; k < len(element.Options); k++ {
				opt := element.Options[k]
				option := forms.Option{
					ElementID: elemId,
					Name:      opt.Name,
					Position:  opt.Index,
				}
				insertOption := fmt.Sprintf("INSERT INTO options (elementID, name, position) VALUES (%v, \"%s\", %v)", option.ElementID, option.Name, option.Position)
				_, err := db.Exec(insertOption)
				if err != nil {
					fmt.Println("Failed SQL: " + insertOption)
					log.Fatal("Failed to insert option: " + err.Error())
				}
			}
		}
		count++
		fmt.Printf("Successfully submitted complete form with ID %v and name %s\n", id, sqlForm.Name)
	}
	fmt.Printf("Successfully migrated %v forms to icc-dev on Planetscale\n", count)
}
