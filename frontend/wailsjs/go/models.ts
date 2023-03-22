export namespace main {
	
	export class ChatProcessReq {
	    prompt: string;
	    // Go type: struct { ConversationId string "json:\"conversationId,omitempty\""; ParentMessageId string "json:\"parentMessageId,omitempty\"" }
	    options?: any;
	
	    static createFrom(source: any = {}) {
	        return new ChatProcessReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.prompt = source["prompt"];
	        this.options = this.convertValues(source["options"], Object);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
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

