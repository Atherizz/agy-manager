package profile

import (

    "encoding/json"

    "errors"

    "os"

)

type State struct {

    ActiveProfile string `json:"active_profile"`

}

func SaveState(p *PathResolver, s *State) error {

    data, err := json.MarshalIndent(s, "", "  ")

    if err != nil {

        return err

    }

    return os.WriteFile(p.StateFile(), data, 0644)

}

func LoadState(p *PathResolver) (*State, error) {

    data, err := os.ReadFile(p.StateFile())

    if err != nil {

        if errors.Is(err, os.ErrNotExist) {

            return &State{}, nil

        }

        return nil, err

    }

    var s State

    if err := json.Unmarshal(data, &s); err != nil {

        return nil, err

    }

    return &s, nil

}
