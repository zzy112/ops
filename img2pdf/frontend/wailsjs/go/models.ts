export namespace model {
	
	export class FileInfo {
	    path: string;
	    name: string;
	    size: number;
	    thumbnail: string;
	
	    static createFrom(source: any = {}) {
	        return new FileInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.size = source["size"];
	        this.thumbnail = source["thumbnail"];
	    }
	}

}

