package db

// func Test_dbRepository_GetCategoryNameDB(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		categoryId int64
// 	}
// 	tests := []struct {
// 		name     string
// 		dbR      *dbRepository
// 		args     args
// 		wantName string
// 		wantErr  bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotName, err := tt.dbR.GetCategoryNameDB(tt.args.ctx, tt.args.categoryId)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("dbRepository.GetCategoryNameDB() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if gotName != tt.wantName {
// 				t.Errorf("dbRepository.GetCategoryNameDB() = %v, want %v", gotName, tt.wantName)
// 			}
// 		})
// 	}
// }

// func Test_dbRepository_GetMonthSumDB(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		categoryId int64
// 	}
// 	tests := []struct {
// 		name       string
// 		dbR        *dbRepository
// 		args       args
// 		wantPrices []int64
// 		wantErr    bool
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			gotPrices, err := tt.dbR.GetMonthSumDB(tt.args.ctx, tt.args.categoryId)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("dbRepository.GetMonthSumDB() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			if !reflect.DeepEqual(gotPrices, tt.wantPrices) {
// 				t.Errorf("dbRepository.GetMonthSumDB() = %v, want %v", gotPrices, tt.wantPrices)
// 			}
// 		})
// 	}
// }

// func Test_dbRepository_AddRecordDB(t *testing.T) {
// 	type args struct {
// 		ctx        context.Context
// 		categoryId int64
// 		price      int64
// 		date       string
// 	}
// 	tests := []struct {
// 		name string
// 		dbR  *dbRepository
// 		args args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			tt.dbR.AddRecordDB(tt.args.ctx, tt.args.categoryId, tt.args.price, tt.args.date)
// 		})
// 	}
// }
