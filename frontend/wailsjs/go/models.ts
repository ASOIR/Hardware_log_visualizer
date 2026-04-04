export namespace main {
	
	export class ParseResult {
	    headers: string[];
	    times: string[];
	    data: number[][];
	    files: string[];
	    loadedFile: string;
	
	    static createFrom(source: any = {}) {
	        return new ParseResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.headers = source["headers"];
	        this.times = source["times"];
	        this.data = source["data"];
	        this.files = source["files"];
	        this.loadedFile = source["loadedFile"];
	    }
	}

}

