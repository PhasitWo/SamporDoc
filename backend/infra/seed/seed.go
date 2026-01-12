package seed

import (
	"SamporDoc/backend/infra/model"
	"SamporDoc/backend/infra/repository"
	"SamporDoc/backend/utils"
	"fmt"
)

var DefaultShops []model.Shop = []model.Shop{
	{Slug: "HJKmain", Name: "หจก (หลัก)", SortingLevel: 0},
	{Slug: "HJKsec", Name: "หจก (รอง)", SortingLevel: 1},
	{Slug: "BUMmain", Name: "บำเหน็จ (หลัก)", SortingLevel: 2},
	{Slug: "BUMsec", Name: "บำเหน็จ (รอง)", SortingLevel: 3},
}

var DefaultCustomers []model.Customer = []model.Customer{
	{
		Name:    "โรงเรียนเขาดินพิทยารักษ์",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนเทพสถิตวิทยา",
		Address: utils.Ptr("อ.เทพสถิต จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนเริงรมย์วิทยาคม",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนโคกเพชรวิทยาคาร",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนโคกหินตั้งศึกษาศิลป์",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนโนนจาน(เนตรขันธ์ราษฎร์บำรุง)",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนคอนสารวิทยาคม",
		Address: utils.Ptr("อ.คอนสาร จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนชุมชนบ้านเพชร(วันครู2500)",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนชุมชนบ้านหนองแวง(คุรุราษฎร์อุปถัมภ์)",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนตาเนินราษฎร์วิทยาคาร",
		Address: utils.Ptr("อ.เนินสง่า จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านแก้งยาว",
		Address: utils.Ptr("อ.บ้านเขว้า จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านโคกแพงพวย(บัวประชาสรรค์)",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านโคกแสว",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านโคกไค",
		Address: utils.Ptr("อ.เทพสถิต จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านโคกสว่าง",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านโคกสะอาด",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านโปร่งมีชัย",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านโสกรวกหนองซึก",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านกลอยสามัคคี",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านดอนละนาม",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านตาล",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านทองคำพิงวิทยา",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านปะโคตะกอ(ประชาสามัคคี)",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านยางเครือ",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านละหาน(อภิรักษ์วิทยา)",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหนองโดน",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหนองดง",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหนองผักแว่น",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหนองลูกช้าง",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหนองหญ้าข้าวนก",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหนองอีหล่อ",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหลุบโพธิ์",
		Address: utils.Ptr("อ.บ้านเขว้า จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหลุบงิ้ว",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านห้วย",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านห้วยยาง",
		Address: utils.Ptr("อ.จัตุรัส จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนบ้านหัวทะเล(ผ่องประชาสรรค์)",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนมะนะศึกษา",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนวังกะอาม",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "โรงเรียนหนองกกสามัคคี",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "ศูนย์พัฒนาเด็กเล็ก อบต.โคกเริงรมย์",
		Address: utils.Ptr("อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "เทศบาลตำบลบำเหน็จณรงค์",
		Address: utils.Ptr("555 ม.11 ต.บ้านชวน อ.บำเหน็จณรงค์ จ.ชัยภูมิ"),
	},
	{
		Name:    "การไฟฟ้าส่วนภูมิภาคอำเภอบำเหน็จณรงค์",
		Address: utils.Ptr("138 หมู่ 5 ถ.รุ่งเรืองศรี ต.บ้านชวน อ.บำเหน็จณรงค์ จ.ชัยภูมิ 36160"),
	},
	{
		Name:    "สำนักงานพัฒนาชุมชน อ.เทพสถิต จ.ชัยภูมิ",
		Address: utils.Ptr("หมู่ 1 ต.วะตะแบก อ.เทพสถิต จ.ชัยภูมิ 36230"),
	},
	{
		Name:    "โรงเรียนภูเขียว",
		Address: utils.Ptr("อ.เภอภูเขียว จ.ชัยภูมิ"),
	},
}

func SeedShops(repo *repository.Repo) {
	fmt.Println("Start seeding shops...")
	for _, shop := range DefaultShops {
		err := repo.CreateShop(&shop)
		if err != nil {
			fmt.Println("Error seeding database -> repo.CreateShop")
			panic(err)
		}
	}
	fmt.Println("Finish seeding shops")
}

func SeedCustomers(repo *repository.Repo) {
	fmt.Println("Start seeding customers...")
	for _, sch := range DefaultCustomers {
		err := repo.CreateCustomer(&sch)
		if err != nil {
			fmt.Println("Error seeding database -> repo.CreateCustomer")
			panic(err)
		}
	}
	fmt.Println("Finish seeding customers")
}
