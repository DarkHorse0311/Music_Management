package database

import (
	"errors"
	"log"
	"os"
	"path/filepath"
)

func (database *Database) InsertFeedback(time int64, feedback string, header string) error {
	_, err := database.stmt.insertFeedback.Exec(time, feedback, header)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) GetRandomFiles(limit int64) ([]File, error) {
	rows, err := database.stmt.getRandomFiles.Query(limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]File, 0)
	for rows.Next() {
		file := File{
			Db: database,
		}
		err = rows.Scan(&file.ID, &file.Folder_id, &file.Filename, &file.Foldername, &file.Filesize)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (database *Database) GetFilesInFolder(folder_id int64, limit int64, offset int64) ([]File, error) {
	rows, err := database.stmt.getFilesInFolder.Query(folder_id, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	files := make([]File, 0)
	for rows.Next() {
		file := File{
			Db:        database,
			Folder_id: folder_id,
		}
		err = rows.Scan(&file.ID, &file.Filename, &file.Filesize, &file.Foldername)
		if err != nil {
			return nil, err
		}
		files = append(files, file)
	}
	return files, nil
}

func (database *Database) SearchFolders(foldername string, limit int64, offset int64) ([]Folder, error) {
	rows, err := database.stmt.searchFolders.Query("%"+foldername+"%", limit, offset)
	if err != nil {
		return nil, errors.New("Error searching folders at query " + err.Error())
	}
	defer rows.Close()
	folders := make([]Folder, 0)
	for rows.Next() {
		folder := Folder{
			Db: database,
		}
		err = rows.Scan(&folder.ID, &folder.Folder, &folder.Foldername)
		if err != nil {
			return nil, errors.New("Error scanning SearchFolders" + err.Error())
		}
		folders = append(folders, folder)
	}
	return folders, nil
}

func (database *Database) GetFile(id int64) (*File, error) {
	file := &File{
		Db: database,
	}
	err := database.stmt.getFile.QueryRow(id).Scan(&file.ID, &file.Folder_id, &file.Filename, &file.Foldername, &file.Filesize)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (database *Database) ResetFiles() error {
	log.Println("[db] Reset files")
	var err error
	_, err = database.stmt.dropFiles.Exec()
	if err != nil {
		return err
	}
	_, err = database.stmt.initFilesTable.Exec()
	if err != nil {
		return err
	}
	return err
}

func (database *Database) ResetFolder() error {
	log.Println("[db] Reset folders")
	var err error
	_, err = database.stmt.dropFolder.Exec()
	if err != nil {
		return err
	}
	_, err = database.stmt.initFoldersTable.Exec()
	if err != nil {
		return err
	}
	return err
}

func (database *Database) Walk(root string, pattern []string) error {
	patternDict := make(map[string]bool)
	for _, v := range pattern {
		patternDict[v] = true
	}
	log.Println("[db] Walk", root, patternDict)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		// check pattern
		ext := filepath.Ext(info.Name())
		if _, ok := patternDict[ext]; !ok {
			return nil
		}

		// insert file, folder will aut created
		err = database.Insert(path, info.Size())
		if err != nil {
			return err
		}
		return nil
	})
}

func (database *Database) GetFolder(folderId int64) (*Folder, error) {
	folder := &Folder{
		Db: database,
	}
	err := database.stmt.getFolder.QueryRow(folderId).Scan(&folder.Folder)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

func (database *Database) SearchFiles(filename string, limit int64, offset int64) ([]File, error) {
	rows, err := database.stmt.searchFiles.Query("%"+filename+"%", limit, offset)
	if err != nil {
		return nil, errors.New("Error searching files at query " + err.Error())
	}
	defer rows.Close()
	files := make([]File, 0)
	for rows.Next() {
		var file File = File{
			Db: database,
		}
		err = rows.Scan(&file.ID, &file.Folder_id, &file.Filename, &file.Foldername, &file.Filesize)
		if err != nil {
			return nil, errors.New("Error scanning SearchFiles " + err.Error())
		}
		files = append(files, file)
	}
	if err = rows.Err(); err != nil {
		return nil, errors.New("Error scanning SearchFiles exit without full result" + err.Error())
	}
	return files, nil
}

func (database *Database) FindFolder(folder string) (int64, error) {
	var id int64
	err := database.stmt.findFolder.QueryRow(folder).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (database *Database) InsertFolder(folder string) (int64, error) {
	result, err := database.stmt.insertFolder.Exec(folder, filepath.Base(folder))
	if err != nil {
		return 0, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lastInsertId, nil
}

func (database *Database) InsertFile(folderId int64, filename string, filesize int64) error {
	_, err := database.stmt.insertFile.Exec(folderId, filename, filesize)
	if err != nil {
		return err
	}
	return nil
}

func (database *Database) Insert(path string, filesize int64) error {
	folder, filename := filepath.Split(path)
	folderId, err := database.FindFolder(folder)
	if err != nil {
		folderId, err = database.InsertFolder(folder)
		if err != nil {
			return err
		}
	}
	err = database.InsertFile(folderId, filename, filesize)
	if err != nil {
		return err
	}
	return nil
}