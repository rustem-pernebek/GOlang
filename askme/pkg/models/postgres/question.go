package postgres

import (
	"askme/pkg/models"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	askQuestion = "insert into questions(user_id,question) values ($1,$2)"
	answer= "insert into answer(q_id,u_id,answer) values ($1,$2,$3)"
	answersToQuestion="select id,answer from answer where q_id=$1"
	getAllQuestions="select id, question from questions"
	getSingleQuestion="select question from questions where id=$1"

	deleteQuestion="delete from questions where id=$1"
	deleteAnswer="delete from answer where id=$1"
	deleteAllAnswer="delete from answer where q_id=$1"
)
type QuestionModel struct {
	Pool *pgxpool.Pool
}

func (m *QuestionModel) AskQuestion(id int,question string) error  {
	row:=m.Pool.QueryRow(context.Background(),askQuestion,id,question)
	err:=row.Scan()
	if err!=nil {
		return err
	}
	return nil
}

func (m *QuestionModel) DeleteQuestion(id int) error  {
	row:=m.Pool.QueryRow(context.Background(),deleteQuestion,id)
	err:=row.Scan()
	if err!=nil {
		return err
	}
	err=m.DeleteAllAnswers(id)
	if err!=nil {
		return err
	}
	return nil
}
func (m *QuestionModel) DeleteAllAnswers(id int) error  {
	row:=m.Pool.QueryRow(context.Background(),deleteAllAnswer,id)
	err:=row.Scan()
	if err!=nil {
		return err
	}
	return nil
}

func (m *QuestionModel) GetAllQuestions() ([]*models.Question,error)  {

	qs :=[]*models.Question{}

	rows, err :=m.Pool.Query(context.Background(),getAllQuestions)

	if err != nil {
		return nil, err
	}

	for rows.Next(){
		
		q := &models.Question{}
		err:=rows.Scan(&q.Id,&q.Quest)

		if err != nil {

			return nil, err
		}
		q.Answers,_=m.GetAnswers(q.Id)
		qs= append(qs, q)



	}


	if err = rows.Err(); err != nil {
		return nil, err
	}

	return qs,nil
}
func (m *QuestionModel) GetSingleQuestion(id int) (*models.Question,error)  {

	q:=&models.Question{}
	row:=m.Pool.QueryRow(context.Background(),getSingleQuestion,id)

	err:=row.Scan(&q.Quest)
	q.Id=id
	anss,err:=m.GetAnswers(id)
	fmt.Println(q.Quest)
	q.Answers=anss
	if anss==nil {
		return q,nil
	}
	if err != nil {
		return nil, err
	}

	return q,nil
}



func (m *QuestionModel) AnswerQuestion(id,q_id int,ans string) error  {
	row:=m.Pool.QueryRow(context.Background(),answer,q_id,id,ans)
	err:=row.Scan()
	if err!=nil {
		return err
	}
	return nil
}
func (m *QuestionModel) DeleteAnswer(id int) error  {
	row:=m.Pool.QueryRow(context.Background(),deleteAnswer,id)
	err:=row.Scan()
	if err!=nil {
		return err
	}
	return nil
}
func (m *QuestionModel)GetAnswers(q_id int) ([]*models.Answer,error) {
	ans :=[]*models.Answer{}
	rows, err :=m.Pool.Query(context.Background(),answersToQuestion,q_id)
	if err != nil {
		return nil, err
	}


	for rows.Next(){

		a:= &models.Answer{}
		err:=rows.Scan(&a.Id,&a.Answer)

		if err != nil {

			return nil, err
		}
		ans= append(ans,a)


	}


	if err = rows.Err(); err != nil {
		return nil, err
	}
	return ans,nil

}


