package rest

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/iDevoid/bitsum/internal/constant/model"
	"github.com/iDevoid/bitsum/internal/module/coins"
	coins_mock "github.com/iDevoid/bitsum/mocks/coins"
	"github.com/valyala/fasthttp"
)

func TestCoinsInit(t *testing.T) {
	type args struct {
		coinsCase coins.Usecase
	}
	tests := []struct {
		name string
		args args
		want CoinsHandler
	}{
		{
			name: "success",
			args: args{
				coinsCase: &coins_mock.MockUsecase{},
			},
			want: &coinsHandler{
				coinsCase: &coins_mock.MockUsecase{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CoinsInit(tt.args.coinsCase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CoinsInit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_coinsHandler_Payment(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := fiber.New()
	type fields struct {
		coinsCase coins.Usecase
	}
	type args struct {
		ctx *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error payload",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "bad date",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"amount": 10
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "bad amount",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"datetime": "2019-10-05T14:48:01+01:00",
						"amount": 0
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "bad amount negative",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"datetime": "2019-10-05T14:48:01+01:00",
						"amount": -1
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "error",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"datetime": "2021-04-12T10:55:23+07:00",
						"amount": 100
					}`))
					return ctx
				}(),
			},
			fields: fields{
				coinsCase: func() coins.Usecase {
					mocked := coins_mock.NewMockUsecase(ctrl)
					mocked.EXPECT().Pay(gomock.Any(), gomock.Any()).Return(errors.New("ERROR"))
					return mocked
				}(),
			},
		},
		{
			name: "success",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"datetime": "2021-04-12T10:55:23+07:00",
						"amount": 100
					}`))
					return ctx
				}(),
			},
			fields: fields{
				coinsCase: func() coins.Usecase {
					mocked := coins_mock.NewMockUsecase(ctrl)
					mocked.EXPECT().Pay(gomock.Any(), gomock.Any()).Return(nil)
					return mocked
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &coinsHandler{
				coinsCase: tt.fields.coinsCase,
			}
			if err := ch.Payment(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("coinsHandler.Payment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_coinsHandler_Receive(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := fiber.New()
	type fields struct {
		coinsCase coins.Usecase
	}
	type args struct {
		ctx *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error payload",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "bad date",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"amount": 10
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "bad amount",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"datetime": "2019-10-05T14:48:01+01:00",
						"amount": 0
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "bad amount negative",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"datetime": "2019-10-05T14:48:01+01:00",
						"amount": -1
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "error",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"datetime": "2021-04-12T10:55:23+07:00",
						"amount": 100
					}`))
					return ctx
				}(),
			},
			fields: fields{
				coinsCase: func() coins.Usecase {
					mocked := coins_mock.NewMockUsecase(ctrl)
					mocked.EXPECT().Receive(gomock.Any(), gomock.Any()).Return(errors.New("ERROR"))
					return mocked
				}(),
			},
		},
		{
			name: "success",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"datetime": "2021-04-12T10:55:23+07:00",
						"amount": 100
					}`))
					return ctx
				}(),
			},
			fields: fields{
				coinsCase: func() coins.Usecase {
					mocked := coins_mock.NewMockUsecase(ctrl)
					mocked.EXPECT().Receive(gomock.Any(), gomock.Any()).Return(nil)
					return mocked
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &coinsHandler{
				coinsCase: tt.fields.coinsCase,
			}
			if err := ch.Receive(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("coinsHandler.Receive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_coinsHandler_History(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := fiber.New()
	type fields struct {
		coinsCase coins.Usecase
	}
	type args struct {
		ctx *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error payload",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "bad start date",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"start_date_time": 
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "bad end date",
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"start_date_time": "2019-10-05T14:48:01+01:00"
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "error",
			fields: fields{
				coinsCase: func() coins.Usecase {
					mocked := coins_mock.NewMockUsecase(ctrl)
					mocked.EXPECT().HistoryTransaction(gomock.Any(), gomock.Any()).Return([]model.Coin{}, errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"start_date_time": "2019-10-05T14:48:01+01:00",
						"end_date_time": "2019-10-06T14:48:01+01:00"
					}`))
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				coinsCase: func() coins.Usecase {
					mocked := coins_mock.NewMockUsecase(ctrl)
					mocked.EXPECT().HistoryTransaction(gomock.Any(), gomock.Any()).Return([]model.Coin{
						{
							DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
							Amount:   10,
						},
					}, nil)
					return mocked
				}(),
			},
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					ctx.Request().SetBody([]byte(`{
						"start_date_time": "2019-10-05T14:48:01+01:00",
						"end_date_time": "2019-10-06T14:48:01+01:00"
					}`))
					return ctx
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &coinsHandler{
				coinsCase: tt.fields.coinsCase,
			}
			if err := ch.History(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("coinsHandler.History() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_coinsHandler_Wallet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	app := fiber.New()
	type fields struct {
		coinsCase coins.Usecase
	}
	type args struct {
		ctx *fiber.Ctx
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "error",
			fields: fields{
				coinsCase: func() coins.Usecase {
					mocked := coins_mock.NewMockUsecase(ctrl)
					mocked.EXPECT().Balance(gomock.Any()).Return(model.Coin{}, errors.New("ERROR"))
					return mocked
				}(),
			},
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					return ctx
				}(),
			},
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				coinsCase: func() coins.Usecase {
					mocked := coins_mock.NewMockUsecase(ctrl)
					mocked.EXPECT().Balance(gomock.Any()).Return(model.Coin{
						DateTime: time.Date(2021, time.April, 12, 10, 55, 23, 0, time.UTC),
						Amount:   10,
					}, nil)
					return mocked
				}(),
			},
			args: args{
				ctx: func() *fiber.Ctx {
					ctx := app.AcquireCtx(&fasthttp.RequestCtx{})
					return ctx
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := &coinsHandler{
				coinsCase: tt.fields.coinsCase,
			}
			if err := ch.Wallet(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("coinsHandler.Wallet() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
