package main

import "fmt"

func getMontlyArchive(user string, year, month int) (*monthDataArchive, error) {
	doc := &monthDataArchive{
		ID:        fmt.Sprintf("%s_%d-%d", user, year, month),
		Positions: []position{},
	}

	if err := couchConn.Get(doc.ID, doc); err != nil {
		return doc, err
	}

	return doc, nil
}

func saveMonthlyArchive(archive *monthDataArchive) error {
	res, err := couchConn.Save(archive)

	if err != nil {
		return err
	}

	if !res.Ok {
		return fmt.Errorf("Could not store archive: %s / %s", res.Reason, res.Error)
	}

	archive.Rev = res.Rev

	return nil
}
