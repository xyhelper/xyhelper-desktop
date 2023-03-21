export namespace main {
	
	export class SessionRes {
	    auth: boolean;
	    model: string;
	
	    static createFrom(source: any = {}) {
	        return new SessionRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.auth = source["auth"];
	        this.model = source["model"];
	    }
	}

}

