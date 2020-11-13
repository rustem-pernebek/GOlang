package postgres

import (
	"askme/pkg/models"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	insertPerson  ="insert into users(login,password,name,phone) values ($1,$2,$3,$4)"
	insertAdmin  ="insert into admin(pers_id) values ($1)"
	getPerson="select*from users where login = $1 and password = $2"
	getAllPerson="select*from person"
	getPersonByID= "select*from users where login = $1"
	getAdmin = "select  u.id, login, password, name, phone from  users u,admin a where a.u_id =u.id and login=$1 and password=$2"

)
type PersonModel struct {
	Pool *pgxpool.Pool
}
func (m *PersonModel) SignUp(person *models.User) (error) {

	row :=m.Pool.QueryRow(context.Background(),insertPerson,person.Login,person.Passrword,person.Name,person.Phone)
	err := row.Scan()
	if err!=nil {
		return err
	}

	return nil
}


func (m *PersonModel) GetPerson(login string,pass string) (*models.User,error) {
	p:=&models.User{}

	err :=m.Pool.QueryRow(context.Background(),getPerson,login,pass).
		Scan(&p.Id,&p.Login,&p.Passrword,&p.Name,&p.Phone)
	if err != nil {
		return nil, err
	}

	return p,nil

}
func (m *PersonModel) GetAllPersons(login string,pass string) ([]*models.User,error) {
	ps:=[]*models.User{}

	rows,err :=m.Pool.Query(context.Background(),getAllPerson,login,pass)
	if err != nil {
		return nil, err
	}
	for rows.Next(){
		p:=&models.User{}
		rows.Scan(&p.Id,&p.Login,&p.Passrword,&p.Name,&p.Phone)
		ps = append(ps, p)
	}


	return ps,nil

}

func (m *PersonModel) GetAdmin(login string,pass string) (*models.User,error) {
	p:=&models.User{}

	err :=m.Pool.QueryRow(context.Background(),getAdmin,login,pass).
		Scan(&p.Id,&p.Login,&p.Passrword,&p.Name,&p.Phone)
	if err != nil {

		return nil, err
	}

	return p,nil

}

func (m *PersonModel) GetPersonBylogin(id string)(*models.User,error){
	p:=&models.User{}


	err :=m.Pool.QueryRow(context.Background(),getPersonByID,id).
		Scan(&p.Id,&p.Login,&p.Passrword,&p.Name,&p.Phone)
	if err != nil {

		return nil, err
	}

	return p,nil
}


