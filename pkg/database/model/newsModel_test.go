package model

import (
	"newsaggr/pkg/config"
	"reflect"
	"testing"
)

// preload - загрузка конфигурации
func preload() error {
	if err := config.Init(); err != nil {
		return err
	}

	return nil
}

func TestDelete(t *testing.T) {
	if err := preload(); err != nil {
		t.Errorf("%v", err)
	}

	type args struct {
		id uint
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Удаление записи",
			args:    args{id: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFindAll(t *testing.T) {
	if err := preload(); err != nil {
		t.Errorf("%v", err)
	}

	tests := []struct {
		name    string
		want    []*News
		wantErr bool
	}{
		{
			name: "Поиск всех записей",
			want: []*News{{
				ID:      1,
				Title:   "Title",
				Content: "Content",
				PubTime: 0,
				Link:    "localhost",
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindAll() got = %v, want %v", got, tt.want)
			}

			news := make([]News, len(got))
			for i, row := range got {
				news[i] = *row
			}
			t.Logf("Все записи: %+v", news)
		})
	}
}

func TestNews_Create(t *testing.T) {

	if err := preload(); err != nil {
		t.Errorf("%v", err)
	}

	type fields struct {
		ID      int
		Title   string
		Content string
		PubTime int64
		Link    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Создание записи",
			fields: fields{
				Title:   "Title",
				Content: "Content",
				PubTime: 0,
				Link:    "localhost",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &News{
				ID:      tt.fields.ID,
				Title:   tt.fields.Title,
				Content: tt.fields.Content,
				PubTime: tt.fields.PubTime,
				Link:    tt.fields.Link,
			}
			if err := n.Create(); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNews_FindOne(t *testing.T) {
	if err := preload(); err != nil {
		t.Errorf("%v", err)
	}

	type fields struct {
		ID      int
		Title   string
		Content string
		PubTime int64
		Link    string
	}
	type args struct {
		id int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Поиск одной записи",
			fields:  fields{},
			args:    args{id: 1},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &News{
				ID:      tt.fields.ID,
				Title:   tt.fields.Title,
				Content: tt.fields.Content,
				PubTime: tt.fields.PubTime,
				Link:    tt.fields.Link,
			}
			if err := n.FindOne(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("FindOne() error = %v, wantErr %v", err, tt.wantErr)
			}
			t.Logf("Новость: %v", *n)
		})
	}
}

func TestFindLimit(t *testing.T) {
	if err := preload(); err != nil {
		t.Errorf("%v", err)
	}

	type args struct {
		limit int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Определенное количество записей",
			args:    args{limit: 5},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := FindLimit(tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindLimit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Новости: %+v", got)
		})
	}
}