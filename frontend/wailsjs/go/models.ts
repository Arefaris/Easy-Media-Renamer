export namespace config {
	
	export class Config {
	    tmdb_api_key: string;
	    name_template: string;
	    media_extensions: string[];
	    include_specials: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tmdb_api_key = source["tmdb_api_key"];
	        this.name_template = source["name_template"];
	        this.media_extensions = source["media_extensions"];
	        this.include_specials = source["include_specials"];
	    }
	}

}

export namespace media {
	
	export class File {
	    name: string;
	    path: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.size = source["size"];
	    }
	}
	export class Op {
	    from: string;
	    to: string;
	    status: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Op(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = source["from"];
	        this.to = source["to"];
	        this.status = source["status"];
	        this.error = source["error"];
	    }
	}
	export class Result {
	    renamed: number;
	    skipped: number;
	    failed: Op[];
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.renamed = source["renamed"];
	        this.skipped = source["skipped"];
	        this.failed = this.convertValues(source["failed"], Op);
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

export namespace provider {
	
	export class Episode {
	    season: number;
	    number: number;
	    title: string;
	    special: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Episode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.season = source["season"];
	        this.number = source["number"];
	        this.title = source["title"];
	        this.special = source["special"];
	    }
	}
	export class Info {
	    name: string;
	    configured: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.configured = source["configured"];
	    }
	}
	export class Show {
	    id: string;
	    name: string;
	    year: string;
	    kind: string;
	    genres: string[];
	
	    static createFrom(source: any = {}) {
	        return new Show(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.year = source["year"];
	        this.kind = source["kind"];
	        this.genres = source["genres"];
	    }
	}

}

