package exec

import "fmt"

func (p *preprocessor) preprocessGrantOrRevokeKey(patternKey string, params ...interface{}) (int, error) {
	params = append(params, p.grantOrRevoke)
	path := fmt.Sprintf(pathPatterns[patternKey], params...)
	count, err := p.makeFileInstructions(path)
	if err != nil {
		return count, err
	}
	return count, nil
}


