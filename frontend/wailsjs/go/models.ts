export namespace main {
	
	export class CreateReceiptParams {
	    TemplatePath: string;
	    ControlPath: string;
	    Filename: string;
	    OutputDir: string;
	    ReceiptNO: string;
	    ReceiptDate?: string;
	    CustomerName: string;
	    CustomerID?: number;
	    Address?: string;
	    Detail?: string;
	    DeliveryNO?: string;
	    DeliveryDate?: string;
	    Amount: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateReceiptParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TemplatePath = source["TemplatePath"];
	        this.ControlPath = source["ControlPath"];
	        this.Filename = source["Filename"];
	        this.OutputDir = source["OutputDir"];
	        this.ReceiptNO = source["ReceiptNO"];
	        this.ReceiptDate = source["ReceiptDate"];
	        this.CustomerName = source["CustomerName"];
	        this.CustomerID = source["CustomerID"];
	        this.Address = source["Address"];
	        this.Detail = source["Detail"];
	        this.DeliveryNO = source["DeliveryNO"];
	        this.DeliveryDate = source["DeliveryDate"];
	        this.Amount = source["Amount"];
	    }
	}

}

export namespace model {
	
	export class Customer {
	    ID: number;
	    name: string;
	    address?: string | null;
	    headCheckerName?: string | null;
	    checker1Name?: string | null;
	    checker2Name?: string | null;
	    objectName?: string | null;
	    headObjectName?: string | null;
	    bossName?: string | null;
	
	    static createFrom(source: any = {}) {
	        return new Customer(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.headCheckerName = source["headCheckerName"];
	        this.checker1Name = source["checker1Name"];
	        this.checker2Name = source["checker2Name"];
	        this.objectName = source["objectName"];
	        this.headObjectName = source["headObjectName"];
	        this.bossName = source["bossName"];
	    }
	}
	export class Shop {
	    ID: number;
	    slug: string;
	    name: string;
	    sortingLevel: number;
	    receiptFormPath?: string | null;
	    receiptControlPath?: string | null;
	
	    static createFrom(source: any = {}) {
	        return new Shop(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.slug = source["slug"];
	        this.name = source["name"];
	        this.sortingLevel = source["sortingLevel"];
	        this.receiptFormPath = source["receiptFormPath"];
	        this.receiptControlPath = source["receiptControlPath"];
	    }
	}

}

