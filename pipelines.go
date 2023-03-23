package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
)

type Pipelines struct {
	c *Client
}

func (p *Pipelines) List(po *PipelinesOptions) (interface{}, error) {
	urlStr := p.c.requestUrl("/repositories/%s/%s/pipelines/", po.Owner, po.RepoSlug)

	if po.Query != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("q", po.Query)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Sort != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("sort", po.Sort)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Page != 0 {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("page", fmt.Sprint(po.Page))
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	return p.c.executePaginated("GET", urlStr, "")
}

func (p *Pipelines) Get(po *PipelinesOptions) (interface{}, error) {
	urlStr := p.c.requestUrl("/repositories/%s/%s/pipelines/%s", po.Owner, po.RepoSlug, po.IDOrUuid)
	return p.c.execute("GET", urlStr, "")
}

func (p *Pipelines) ListSteps(po *PipelinesOptions) (interface{}, error) {
	urlStr := p.c.requestUrl("/repositories/%s/%s/pipelines/%s/steps/", po.Owner, po.RepoSlug, po.IDOrUuid)

	if po.Query != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("q", po.Query)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Sort != "" {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("sort", po.Sort)
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	if po.Page != 0 {
		parsed, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		query := parsed.Query()
		query.Set("page", fmt.Sprint(po.Page))
		parsed.RawQuery = query.Encode()
		urlStr = parsed.String()
	}

	return p.c.executePaginated("GET", urlStr, "")
}

func (p *Pipelines) GetStep(po *PipelinesOptions) (interface{}, error) {
	urlStr := p.c.requestUrl("/repositories/%s/%s/pipelines/%s/steps/%s", po.Owner, po.RepoSlug, po.IDOrUuid, po.StepUuid)
	return p.c.execute("GET", urlStr, "")
}

func (p *Pipelines) GetLog(po *PipelinesOptions) (string, error) {
	urlStr := p.c.requestUrl("/repositories/%s/%s/pipelines/%s/steps/%s/log", po.Owner, po.RepoSlug, po.IDOrUuid, po.StepUuid)
	responseBody, err := p.c.executeRaw("GET", urlStr, "")
	if err != nil {
		return "", err
	}
	defer responseBody.Close()

	rawBody, err := io.ReadAll(responseBody)
	if err != nil {
		return "", err
	}

	return string(rawBody), nil
}

type BitbucketTrigerPipelineRequestBody struct {
	Target struct {
		RefType  string `json:"ref_type"`
		Type     string `json:"type"`
		RefName  string `json:"ref_name"`
		Selector struct {
			Type    string `json:"type"`
			Pattern string `json:"pattern"`
		} `json:"selector"`
	} `json:"target"`
	Variables []struct {
		Key     string `json:"key"`
		Value   string `json:"value"`
		Secured bool   `json:"secured,omitempty"`
	} `json:"variables"`
}

func (p *Pipelines) TriggerPipeline(po *PipelinesOptions, body *BitbucketTrigerPipelineRequestBody) (interface{}, error) {
	urlStr := p.c.requestUrl("/repositories/%s/%s/pipelines/", po.Owner, po.RepoSlug)

	b, err := json.Marshal(body)
	if err != nil {
		return "failed to parse body", err
	}
	data := string(b)

	responseBody, err := p.c.execute("POST", urlStr, data)
	if err != nil {
		return "failed to trigger bitbucket pipeline", err
	}

	return responseBody, nil
}
