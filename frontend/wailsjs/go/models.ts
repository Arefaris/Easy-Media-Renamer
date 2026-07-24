export namespace config {
	
	export class Config {
	    tmdb_api_key: string;
	    name_template: string;
	    media_extensions: string[];
	    include_specials: boolean;
	    recursive: boolean;
	    max_depth: number;
	    action: string;
	    auto_match: boolean;
	    preset: string;
	    language: string;
	
	    static createFrom(source: any = {}) {
	        return new Config(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.tmdb_api_key = source["tmdb_api_key"];
	        this.name_template = source["name_template"];
	        this.media_extensions = source["media_extensions"];
	        this.include_specials = source["include_specials"];
	        this.recursive = source["recursive"];
	        this.max_depth = source["max_depth"];
	        this.action = source["action"];
	        this.auto_match = source["auto_match"];
	        this.preset = source["preset"];
	        this.language = source["language"];
	    }
	}

}

export namespace match {
	
	export class Pair {
	    file_index: number;
	    episode_index: number;
	    confidence: number;
	    strategy: string;
	
	    static createFrom(source: any = {}) {
	        return new Pair(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file_index = source["file_index"];
	        this.episode_index = source["episode_index"];
	        this.confidence = source["confidence"];
	        this.strategy = source["strategy"];
	    }
	}

}

export namespace media {
	
	export class File {
	    name: string;
	    path: string;
	    rel_path: string;
	    size: number;
	
	    static createFrom(source: any = {}) {
	        return new File(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.rel_path = source["rel_path"];
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
	export class HistoryEntry {
	    id: string;
	    // Go type: time
	    time: any;
	    directory: string;
	    action: string;
	    count: number;
	    ops: Op[];
	
	    static createFrom(source: any = {}) {
	        return new HistoryEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.time = this.convertValues(source["time"], null);
	        this.directory = source["directory"];
	        this.action = source["action"];
	        this.count = source["count"];
	        this.ops = this.convertValues(source["ops"], Op);
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
	
	export class Preset {
	    name: string;
	    template: string;
	
	    static createFrom(source: any = {}) {
	        return new Preset(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.template = source["template"];
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
	    airdate: string;
	    absolute: number;
	
	    static createFrom(source: any = {}) {
	        return new Episode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.season = source["season"];
	        this.number = source["number"];
	        this.title = source["title"];
	        this.special = source["special"];
	        this.airdate = source["airdate"];
	        this.absolute = source["absolute"];
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

export namespace verify {
	
	export class Result {
	    file: string;
	    expected: string;
	    actual: string;
	    ok: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new Result(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.file = source["file"];
	        this.expected = source["expected"];
	        this.actual = source["actual"];
	        this.ok = source["ok"];
	        this.error = source["error"];
	    }
	}

}

