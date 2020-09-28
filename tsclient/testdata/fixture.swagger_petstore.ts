// SwaggerPetstore is a client package to interact with the api.
export class SwaggerPetstore {
    baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    // listPets get /pets
    listPets(body: any, limit: integer) : Promise<Response> {
        let url = '/pets';

        let headers = new(Headers);
        headers.set('Content-Type', 'application/json');

        const params = {
            method: 'GET',
            headers: headers,
            body = JSON.stringify(body),
        };

        return fetch(new Request(this.baseUrl + url, params));
    }

    // createPets post /pets
    createPets() : Promise<Response> {
        let url = '/pets';

        let headers = new(Headers);
        headers.set('Content-Type', 'application/json');

        const params = {
            method: 'POST',
            headers: headers,
        };

        return fetch(new Request(this.baseUrl + url, params));
    }

    // showPetById get /pets/{petId}
    showPetById(petId: string) : Promise<Response> {
        let url = '/pets/{petId}';
        url = url.replace('{petId}', petId.toString();

        let headers = new(Headers);
        headers.set('Content-Type', 'application/json');

        const params = {
            method: 'GET',
            headers: headers,
        };

        return fetch(new Request(this.baseUrl + url, params));
    }
}
