package ws

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/ARUMANDESU/uniclubs-comments-service/pkg/logger"
	"github.com/centrifugal/centrifuge"
)

func TestManager_routeEvent(t *testing.T) {
	type fields struct {
		log            *slog.Logger
		node           *centrifuge.Node
		handlers       map[EventType]EventHandler
		commentService CommentService
	}
	type args struct {
		msg clientMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    centrifuge.PublishReply
		wantErr bool
	}{
		{
			name: "unsupported event",
			fields: fields{
				log:            logger.Plug(),
				node:           nil,
				handlers:       nil,
				commentService: nil,
			},
			args: args{
				msg: clientMessage{
					Event: Event{
						Type: "unsupported",
					},
				},
			},
			want:    centrifuge.PublishReply{},
			wantErr: true,
		},
		{
			name: "create comment",
			fields: fields{
				log:  logger.Plug(),
				node: nil,
				handlers: map[EventType]EventHandler{"create_comment": func(msg clientMessage) (centrifuge.PublishReply, error) {
					return centrifuge.PublishReply{}, nil
				}},
				commentService: nil,
			},
			args: args{
				msg: clientMessage{
					Event: Event{
						Type: "create_comment",
					},
				},
			},
			want:    centrifuge.PublishReply{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Manager{
				log:            tt.fields.log,
				node:           tt.fields.node,
				handlers:       tt.fields.handlers,
				commentService: tt.fields.commentService,
			}
			got, err := m.routeEvent(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Manager.routeEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Manager.routeEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
