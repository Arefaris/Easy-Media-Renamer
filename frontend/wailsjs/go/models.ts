export namespace main {
	
	export class Show {
	    score: number;
	    // Go type: struct { ID int "json:\"id\""; Name string "json:\"name\""; Genre []string "json:\"genres\"" }
	    show: any;
	
	    static createFrom(source: any = {}) {
	        return new Show(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.score = source["score"];
	        this.show = this.convertValues(source["show"], Object);
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

