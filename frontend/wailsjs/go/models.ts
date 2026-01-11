export namespace excel {
	
	export class ControlData {
	    NO: number;
	
	    static createFrom(source: any = {}) {
	        return new ControlData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.NO = source["NO"];
	    }
	}

}

export namespace main {
	
	export class ReceiptData {
	    ReceiptNO: number;
	    // Go type: time
	    ReceiptDate?: any;
	    CustomerName: string;
	    Address?: string;
	    Detail?: string;
	    DeliveryNO?: number;
	    // Go type: time
	    DeliveryDate?: any;
	    Amount: number;
	
	    static createFrom(source: any = {}) {
	        return new ReceiptData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ReceiptNO = source["ReceiptNO"];
	        this.ReceiptDate = this.convertValues(source["ReceiptDate"], null);
	        this.CustomerName = source["CustomerName"];
	        this.Address = source["Address"];
	        this.Detail = source["Detail"];
	        this.DeliveryNO = source["DeliveryNO"];
	        this.DeliveryDate = this.convertValues(source["DeliveryDate"], null);
	        this.Amount = source["Amount"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class CreateReceiptParams {
	    TemplatePath: string;
	    OutputFilePath: string;
	    Data: ReceiptData;
	    ControlPath: string;
	    ControlData: excel.ControlData;
	
	    static createFrom(source: any = {}) {
	        return new CreateReceiptParams(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TemplatePath = source["TemplatePath"];
	        this.OutputFilePath = source["OutputFilePath"];
	        this.Data = this.convertValues(source["Data"], ReceiptData);
	        this.ControlPath = source["ControlPath"];
	        this.ControlData = this.convertValues(source["ControlData"], excel.ControlData);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

export namespace model {
	
	export class School {
	    ID: number;
	    Name: string;
	    Address?: string;
	
	    static createFrom(source: any = {}) {
	        return new School(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Name = source["Name"];
	        this.Address = source["Address"];
	    }
	}
	export class Shop {
	    ID: number;
	    Slug: string;
	    Name: string;
	    BillFormPath?: string;
	    BillcontrolPath?: string;
	
	    static createFrom(source: any = {}) {
	        return new Shop(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ID = source["ID"];
	        this.Slug = source["Slug"];
	        this.Name = source["Name"];
	        this.BillFormPath = source["BillFormPath"];
	        this.BillcontrolPath = source["BillcontrolPath"];
	    }
	}

}

