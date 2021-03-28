// SwaggerPetstore is a client package to interact with the api.
export default class SwaggerPetstore {
    baseUrl: string;

    constructor(baseUrl: string) {
        if (baseUrl === null) {
            this.baseUrl = 'http://petstore.swagger.io/v1';
        } else {
            this.baseUrl = baseUrl;
        }
    }

    // listPets get /pets
    listPets(body: any, limit: number): Promise<Response> {
        let url = '/pets';

        let query = '';
        if (limit !== undefined) {
            if (query.length != 0) { query += '&' }
            query += 'limit=' + encodeURIComponent(limit.toString());
        }
        if (query.length != 0) {
            query = '?' + query;
        }

        let headers = new Headers();
        headers.set('Content-Type', 'application/json');

        const params = {
            method: 'GET',
            headers: headers,
            body: JSON.stringify(body),
        };

        return fetch(new Request(this.baseUrl + url + query, params));
    }

    // createPets post /pets
    createPets(): Promise<Response> {
        let url = '/pets';

        let headers = new Headers();
        headers.set('Content-Type', 'application/json');

        const params = {
            method: 'POST',
            headers: headers,
        };

        return fetch(new Request(this.baseUrl + url, params));
    }

    // showPetById get /pets/{petId}
    showPetById(petId: string): Promise<Response> {
        let url = '/pets/{petId}';
        url = url.replace('{petId}', petId.toString());

        let headers = new Headers();
        headers.set('Content-Type', 'application/json');

        const params = {
            method: 'GET',
            headers: headers,
        };

        return fetch(new Request(this.baseUrl + url, params));
    }
}
