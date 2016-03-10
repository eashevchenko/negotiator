package negotiator

import (
	"encoding/json"
	"net/http"
	"strings"
)

const defaultJsonContentType = "application/json"

type jsonProcessor struct {
	dense          bool
	prefix, indent string
	contentType    string
}

func NewJSON() ResponseProcessor {
	return &jsonProcessor{true, "", "", defaultJsonContentType}
}

func NewJSONIndent(prefix, index string) ResponseProcessor {
	return &jsonProcessor{false, prefix, index, defaultJsonContentType}
}

func NewJSONIndent2Spaces() ResponseProcessor {
	return NewJSONIndent("", "  ")
}

func (p *jsonProcessor) SetContentType(contentType string) ResponseProcessor {
	p.contentType = contentType
	return p
}

func (*jsonProcessor) CanProcess(mediaRange string) bool {
	return strings.EqualFold(mediaRange, "application/json") ||
		strings.HasPrefix(mediaRange, "application/json-") ||
		strings.HasSuffix(mediaRange, "+json")
}

func (p *jsonProcessor) Process(w http.ResponseWriter, dataModel interface{}) error {
	if dataModel == nil {
		w.WriteHeader(http.StatusNoContent)
		return nil

	} else {
		w.Header().Set("Content-Type", p.contentType)
		if p.dense {
			return json.NewEncoder(w).Encode(dataModel)

		} else {
			js, err := json.MarshalIndent(dataModel, p.prefix, p.indent)

			if err != nil {
				return err
			}

			return writeWithNewline(w, js)
		}
	}
}
