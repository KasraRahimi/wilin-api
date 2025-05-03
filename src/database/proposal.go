package database

import (
	"database/sql"
	"fmt"
)

type ProposalModel struct {
	Id     int
	userId int
	Entry  string
	Pos    string
	Gloss  string
	Notes  string
}

type ProposalDao struct {
	Db *sql.DB
}

func (dao *ProposalDao) CreateProposal(proposal *ProposalModel) (int64, error) {
	query := `INSERT INTO proposals (user_id, entry, pos, gloss, notes) 
		VALUES (?, ?, ?, ?, ?)`
	result, err := dao.Db.Exec(
		query,
		proposal.userId,
		proposal.Entry,
		proposal.Pos,
		proposal.Gloss,
		proposal.Notes,
	)
	if err != nil {
		return 0, fmt.Errorf("CreateProposal, failed to insert proposal: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("CreateProposal, failed to fetch proposal id: %w", err)
	}
	return id, nil
}
