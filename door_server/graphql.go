package main

import (
	"encoding/base64"
	"github.com/fluidmediaproductions/hotel_door"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"reflect"
)

var bytesScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "Bytes",
	Description: "Byte array",
	Serialize: func(value interface{}) interface{} {
		bytes, isOK := value.([]byte)
		if isOK {
			base64.StdEncoding.EncodeToString(bytes)
		}
		return nil
	},
	ParseValue: func(value interface{}) interface{} {
		stringValue, isOK := value.(string)
		if isOK {
			base64.StdEncoding.DecodeString(stringValue)
		}
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return valueAST.Value
		}
		return nil
	},
})

var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, isOK := p.Source.(*User)
				if isOK {
					return user.ID, nil
				}
				return nil, nil
			},
		},
		"login": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"isAdmin": &graphql.Field{
			Type: graphql.Boolean,
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, isOK := p.Source.(*User)
				if isOK {
					return user.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, isOK := p.Source.(*User)
				if isOK {
					return user.UpdatedAt, nil
				}
				return nil, nil
			},
		},
		"deletedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				user, isOK := p.Source.(*User)
				if isOK {
					return user.DeletedAt, nil
				}
				return nil, nil
			},
		},
	},
})

var piType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Pi",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				pi, isOK := p.Source.(*Pi)
				if isOK {
					return pi.ID, nil
				}
				return nil, nil
			},
		},
		"mac": &graphql.Field{
			Type: graphql.String,
		},
		"lastSeen": &graphql.Field{
			Type: graphql.DateTime,
		},
		"online": &graphql.Field{
			Type: graphql.Boolean,
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				pi, isOK := p.Source.(*Pi)
				if isOK {
					return pi.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				pi, isOK := p.Source.(*Pi)
				if isOK {
					return pi.UpdatedAt, nil
				}
				return nil, nil
			},
		},
		"deletedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				pi, isOK := p.Source.(*Pi)
				if isOK {
					return pi.DeletedAt, nil
				}
				return nil, nil
			},
		},
	},
})

var actionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Action",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				action, isOK := p.Source.(*Action)
				if isOK {
					return action.ID, nil
				}
				return nil, nil
			},
		},
		"pi": &graphql.Field{
			Type: piType,
		},
		"piId": &graphql.Field{
			Type: graphql.Int,
		},
		"type": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				action, isOK := p.Source.(*Action)
				if isOK {
					return door_comms.DoorAction_name[int32(action.Type)], nil
				}
				return nil, nil
			},
		},
		"complete": &graphql.Field{
			Type: graphql.Boolean,
		},
		"success": &graphql.Field{
			Type: graphql.Boolean,
		},
		"payload": &graphql.Field{
			Type: bytesScalar,
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				door, isOK := p.Source.(*Action)
				if isOK {
					return door.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				door, isOK := p.Source.(*Action)
				if isOK {
					return door.UpdatedAt, nil
				}
				return nil, nil
			},
		},
		"deletedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				door, isOK := p.Source.(*Action)
				if isOK {
					return door.DeletedAt, nil
				}
				return nil, nil
			},
		},
	},
})

var doorType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Door",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				door, isOK := p.Source.(*Door)
				if isOK {
					return door.ID, nil
				}
				return nil, nil
			},
		},
		"pi": &graphql.Field{
			Type: piType,
		},
		"piId": &graphql.Field{
			Type: graphql.Int,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"createdAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				door, isOK := p.Source.(*Door)
				if isOK {
					return door.CreatedAt, nil
				}
				return nil, nil
			},
		},
		"updatedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				door, isOK := p.Source.(*Door)
				if isOK {
					return door.UpdatedAt, nil
				}
				return nil, nil
			},
		},
		"deletedAt": &graphql.Field{
			Type: graphql.DateTime,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				door, isOK := p.Source.(*Door)
				if isOK {
					return door.DeletedAt, nil
				}
				return nil, nil
			},
		},
	},
})

func paginateSlice(arg interface{}, args map[string]interface{}) []interface{} {
	slice, success := takeSliceArg(arg)
	if !success {
		return nil
	}
	offset, isOk := args["offset"].(int)
	if isOk {
		if offset > len(slice) {
			offset = len(slice)
		}
		slice = slice[offset:]
	}
	first, isOk := args["first"].(int)
	if isOk {
		if first > len(slice) {
			first = len(slice)
		}
		slice = slice[:first]
	}
	return slice
}

func takeSliceArg(arg interface{}) (out []interface{}, ok bool) {
	slice, success := takeArg(arg, reflect.Slice)
	if !success {
		ok = false
		return
	}
	c := slice.Len()
	out = make([]interface{}, c)
	for i := 0; i < c; i++ {
		out[i] = slice.Index(i).Interface()
	}
	return out, true
}

func takeArg(arg interface{}, kind reflect.Kind) (val reflect.Value, ok bool) {
	val = reflect.ValueOf(arg)
	if val.Kind() == kind {
		ok = true
	}
	return
}

func makeAuthWrapper(field *graphql.Object) *graphql.Field {
	return &graphql.Field{
		Type: field,
		Args: graphql.FieldConfigArgument{
			"token": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			tokenString, isOK := params.Args["token"].(string)
			if isOK {
			    claims, err := verifyJWT(tokenString)
				if err != nil {
					return nil, err
				} else {
					return claims.User, nil
				}
			}
			return nil, nil
		},
	}
}

func makeRequireAdminWrapper(resolver graphql.FieldResolveFn) graphql.FieldResolveFn {
	return func(params graphql.ResolveParams) (interface{}, error) {
		user, isOk := params.Source.(*User)
		if isOk {
			if user.IsAdmin {
				resolver(params)
			}
		}
		return nil, nil
	}
}

var authenticatedQueries = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthenticatedQueries",
	Fields: graphql.Fields{
		"self": &graphql.Field{
			Type: userType,
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				user, isOk := params.Source.(*User)
				if isOk {
					return user, nil
				}
				return nil, nil
			},
		},

		"pi": &graphql.Field{
			Type: piType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(int)
				if isOK {
					pi := &Pi{}
					err := db.First(pi, idQuery).Error
					if err != nil {
						return nil, err
					}
					return pi, nil
				}
				return nil, nil
			},
		},

		"piList": &graphql.Field{
			Type: graphql.NewList(piType),
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"offset": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				pis := make([]*Pi, 0)
				err := db.Find(&pis).Error
				if err != nil {
					return nil, err
				}
				paginatedPis := paginateSlice(pis, p.Args)
				return paginatedPis, nil
			},
		},

		"door": &graphql.Field{
			Type: doorType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(int)
				if isOK {
					door := &Door{}
					err := db.Set("gorm:auto_preload", true).First(door, idQuery).Error
					if err != nil {
						return nil, err
					}
					return door, nil
				}
				return nil, nil
			},
		},

		"doorList": &graphql.Field{
			Type: graphql.NewList(doorType),
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"offset": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				doors := make([]*Door, 0)
				err := db.Set("gorm:auto_preload", true).Find(&doors).Error
				if err != nil {
					return nil, err
				}
				paginatedDoors := paginateSlice(doors, p.Args)
				return paginatedDoors, nil
			},
		},

		"action": &graphql.Field{
			Type: actionType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				idQuery, isOK := params.Args["id"].(int)
				if isOK {
					action := &Action{}
					err := db.Set("gorm:auto_preload", true).First(action, idQuery).Error
					if err != nil {
						return nil, err
					}
					return action, nil
				}
				return nil, nil
			},
		},

		"actionList": &graphql.Field{
			Type: graphql.NewList(actionType),
			Args: graphql.FieldConfigArgument{
				"first": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"offset": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				actions := make([]*Action, 0)
				err := db.Set("gorm:auto_preload", true).Order("id desc").Find(&actions).Error
				if err != nil {
					return nil, err
				}
				paginatedActions := paginateSlice(actions, p.Args)
				return paginatedActions, nil
			},
		},
	},
})

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"auth": makeAuthWrapper(authenticatedQueries),
	},
})

var authenticatedMutations = graphql.NewObject(graphql.ObjectConfig{
	Name: "AuthenticatedMutations",
	Fields: graphql.Fields{
		"updateDoor": &graphql.Field{
			Type: doorType,
			Args: graphql.FieldConfigArgument{
				"piId": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: makeRequireAdminWrapper(func(params graphql.ResolveParams) (interface{}, error) {
				id, isOK := params.Args["id"].(int)
				if isOK {
					door := &Door{}
					err := db.First(door, id).Error
					if err != nil {
						return nil, err
					}

					piId, isOK := params.Args["piId"].(int)
					if isOK {
						door.PiID = uint(piId)

						db.Model(&Door{}).Where(&Door{PiID: uint(piId)}).Update("pi_id", nil)
					}

					err = db.Save(door).Error
					if err != nil {
						return nil, err
					}
					err = db.Set("gorm:auto_preload", true).First(door).Error
					if err != nil {
						return nil, err
					}

					return door, nil
				}
				return nil, nil
			}),
		},

		"openDoor": &graphql.Field{
			Type: graphql.NewNonNull(graphql.Boolean),
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, isOK := params.Args["id"].(int)
				if isOK {
					door := &Door{}
					err := db.First(door, id).Error
					if err != nil {
						return false, err
					}

					action := &Action{
						PiID: door.PiID,
						Type: int(door_comms.DoorAction_DOOR_UNLOCK),
					}
					err = db.Create(action).Error
					if err != nil {
						return false, err
					}

					return true, nil
				}
				return false, nil
			},
		},

		"deletePi": &graphql.Field{
			Type: piType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: makeRequireAdminWrapper(func(params graphql.ResolveParams) (interface{}, error) {
				id, isOK := params.Args["id"].(int)
				if isOK {
					pi := &Pi{}
					err := db.First(pi, id).Error
					if err != nil {
						return nil, err
					}
					err = db.Delete(pi).Error
					if err != nil {
						return nil, err
					}
					err = db.Unscoped().First(pi).Error
					if err != nil {
						return nil, err
					}

					return pi, nil
				}
				return nil, nil
			}),
		},
	},
})

var rootMutation = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootMutation",
	Fields: graphql.Fields{
		"loginUser": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"login": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"pass": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				login, isOK := params.Args["login"].(string)
				if isOK {
					pass, isOK := params.Args["pass"].(string)
					if isOK {
						user, validLogin := loginUser(login, pass)
						if validLogin {
							tokenString, err := newJWT(user)
							if err != nil {
								return nil, err
							}
							return tokenString, nil
						}
					}
				}
				return nil, nil
			},
		},

		"auth": makeAuthWrapper(authenticatedMutations),

		"refreshToken": &graphql.Field{
			Type: graphql.String,
			Args: graphql.FieldConfigArgument{
				"token": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				tokenString, isOK := params.Args["token"].(string)
				if isOK {
					tokenString, err := refreshJWT(tokenString)
					if err != nil {
						return nil, err
					}
					return tokenString, nil
				}
				return nil, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
	Mutation: rootMutation,
})

func initGraphql() {
	piType.AddFieldConfig("door", &graphql.Field{
		Type: doorType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			pi, isOK := p.Source.(*Pi)
			if isOK {
				door := &Door{}
				err := db.Set("gorm:auto_preload", true).Model(pi).Related(door).Error
				if err != nil {
					return nil, err
				}
				return door, nil
			}
			return nil, nil
		},
	})
}
