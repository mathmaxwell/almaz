package link

import "demo/purpleSchool/pkg/db"

type LinkRepository struct {
	DataBase *db.Db
}

type LinkRepositoryDeps struct {
	DataBase *db.Db
}

func NewLinkRepository(dataBase *db.Db) *LinkRepository {
	return &LinkRepository{
		DataBase: dataBase,
	}

}
func (repo *LinkRepository) Create(Link *Link) (*Link, error) {
	result := repo.DataBase.DB.Create(Link)
	if result.Error != nil {
		return nil, result.Error
	}
	return Link, nil
}
func (repo *LinkRepository) GetByHash(hash string) (*Link, error) {
	var link Link
	result := repo.DataBase.DB.First(&link, "hash = ?", hash)
	if result.Error != nil {
		return nil, result.Error
	}
	return &link, nil
}
