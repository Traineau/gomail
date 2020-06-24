package email

//// Repository struct for db connection
//
//func (repository *Repository) GetUserEmailsFromDiffusionList(id int64) (*[]string, error) {
//	row := repository.Conn.QueryRow("SELECT d.id, d.name, d.description FROM diffusion_list d "+
//		"WHERE c.id=(?)", id)
//	var name, description string
//	switch err := row.Scan(&id, &name, &description, ); err {
//	case sql.ErrNoRows:
//		return nil, nil
//	case nil:
//		diffusionList := DiffusionList{
//			ID:              id,
//			Name:            name,
//			Description:     description,
//		}
//		return &diffusionList, nil
//	default:
//		return nil, err
//	}
//}
//
//func (repository *Repository) SaveDiffusionList(diffusionList *DiffusionList) error {
//	stmt, err := repository.Conn.Prepare("INSERT INTO diffusion_list(name, description) VALUES(?,?)")
//	if err != nil {
//		return err
//	}
//
//	res, errExec := stmt.Exec(diffusionList.Name, diffusionList.Description)
//	if errExec != nil {
//		return errExec
//	}
//
//	lastInsertedID, errInsert := res.LastInsertId()
//	if errInsert != nil {
//		return errInsert
//	}
//
//	diffusionList.ID = lastInsertedID
//
//	return nil
//}
