package database

import (
	"database/sql"
	"fmt"
)

type Scanner interface {
	Scan(dest ...any) error
}

type ProposalModel struct {
	Id     int
	UserId int
	Entry  string
	Pos    string
	Gloss  string
	Notes  string
}

type ProposalUsernameModel struct {
	Id       int
	UserId   int
	Username string
	Entry    string
	Pos      string
	Gloss    string
	Notes    string
}

type ProposalDao struct {
	Db *sql.DB
}

func (dao *ProposalDao) scanProposal(row Scanner) (*ProposalModel, error) {
	var proposal ProposalModel
	err := row.Scan(
		&proposal.Id,
		&proposal.UserId,
		&proposal.Entry,
		&proposal.Pos,
		&proposal.Gloss,
		&proposal.Notes,
	)
	if err != nil {
		return nil, err
	}
	return &proposal, nil
}

func (dao *ProposalDao) scanProposalUsername(row Scanner) (*ProposalUsernameModel, error) {
	var proposal ProposalUsernameModel
	err := row.Scan(
		&proposal.Id,
		&proposal.UserId,
		&proposal.Username,
		&proposal.Entry,
		&proposal.Pos,
		&proposal.Gloss,
		&proposal.Notes,
	)
	if err != nil {
		return nil, err
	}
	return &proposal, nil
}

func (dao *ProposalDao) CreateProposal(proposal *ProposalModel) (int64, error) {
	query := `INSERT INTO proposals (user_id, entry, pos, gloss, notes) VALUES (?, ?, ?, ?, ?)`
	result, err := dao.Db.Exec(
		query,
		proposal.UserId,
		proposal.Entry,
		proposal.Pos,
		proposal.Gloss,
		proposal.Notes,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert proposal: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to fetch proposal id: %w", err)
	}
	return id, nil
}

func (dao *ProposalDao) ReadAllProposals() ([]ProposalModel, error) {
	query := `SELECT id, user_id, entry, pos, gloss, notes FROM proposals`
	rows, err := dao.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query proposals: %w", err)
	}
	defer rows.Close()
	var proposals []ProposalModel
	for rows.Next() {
		proposal, err := dao.scanProposal(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proposals: %w", err)
		}
		proposals = append(proposals, *proposal)
	}
	return proposals, nil
}

func (dao *ProposalDao) ReadAllProposalsWithUsername() ([]ProposalUsernameModel, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.entry, p.pos, p.gloss, p.notes
		FROM proposals p
		JOIN users u ON u.id = p.user_id
		`
	rows, err := dao.Db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query proposals: %w", err)
	}
	defer rows.Close()
	var proposals []ProposalUsernameModel
	for rows.Next() {
		proposal, err := dao.scanProposalUsername(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proposals: %w", err)
		}
		proposals = append(proposals, *proposal)
	}
	return proposals, nil
}

func (dao *ProposalDao) ReadProposalsByUserId(userID int) ([]ProposalModel, error) {
	query := `SELECT id, user_id, entry, pos, gloss, notes FROM proposals WHERE user_id = ?`
	rows, err := dao.Db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query proposals: %w", err)
	}
	defer rows.Close()
	var proposals []ProposalModel
	for rows.Next() {
		proposal, err := dao.scanProposal(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan proposals: %w", err)
		}
		proposals = append(proposals, *proposal)
	}
	return proposals, nil
}

func (dao *ProposalDao) ReadProposalById(id int) (*ProposalModel, error) {
	query := `SELECT id, user_id, entry, pos, gloss, notes FROM proposals WHERE id = ?`
	row := dao.Db.QueryRow(query, id)
	proposal, err := dao.scanProposal(row)
	if err != nil {
		return nil, fmt.Errorf("failed to scan proposal: %w", err)
	}
	return proposal, nil
}

func (dao *ProposalDao) Delete(model *ProposalModel) error {
	query := `DELETE FROM proposals WHERE id = ?`
	_, err := dao.Db.Exec(query, model.Id)
	if err != nil {
		return fmt.Errorf("failed to delete proposal: %w", err)
	}
	return nil
}

func (dao *ProposalDao) Update(model *ProposalModel) error {
	query := `UPDATE proposals 
		SET user_id = ?, entry = ?, pos = ?, gloss = ?, notes = ?
		WHERE id = ?`
	_, err := dao.Db.Exec(query, model.UserId, model.Entry, model.Pos, model.Gloss, model.Notes, model.Id)
	if err != nil {
		return fmt.Errorf("failed to update proposal: %w", err)
	}
	return nil
}
