package models

type sessionModel int

var Sessions sessionModel

func (sessionModel) New(userID int64) (string, error) {

	return "", nil
}

func (sessionModel) Check(sessionID string) (int64, error) {
	return 0, nil
}

func (sessionModel) Delete(sessionID string) error {
	return nil
}
