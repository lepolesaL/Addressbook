import { Injectable } from '@angular/core';
import {Http, Response} from '@angular/http';
import { Headers, RequestOptions } from '@angular/http';
import { Observable } from 'rxjs/Observable';
import { of } from 'rxjs/observable/of';
import 'rxjs/add/operator/map'
import { CONTACTS } from './shared/mockdata'
import { Contact, ApiError } from './shared/contact';
import { SERVER_API_URL } from './app.constants';


@Injectable()
export class AddressbookApiClientService {

  constructor(private http: Http) { }

  getContact(): Observable<any> { 
    return null;
  }

  getContacts(): Observable<any> { 
    return this.http.get('http://localhost:8080/addressbook/contacts').map((res: Response) => { return <Contact[]> res.json()});
    //return of(CONTACTS);
    //const contacts:Contact[] = [];
    //return of(contacts)
  }

  addContact(contact: any): Observable<any> {
    console.log(SERVER_API_URL);
    let body = JSON.stringify(contact);            
    let headers = new Headers({ 'Content-Type': 'application/json', 'Access-Control-Allow-Origin': '*' });
    let options = new RequestOptions({ headers: headers });
    return this.http.post(SERVER_API_URL+'/addressbook/contact', body, options).map(data => {return <ApiError> data.json()});
  }

  editContact(contact: Contact): Observable<any> {
    let body = JSON.stringify(contact);            
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({ headers: headers });
    return this.http.put('http://localhost:8080/addressbook/contact/'+contact.email, body, options).map(data => {return <ApiError> data.json()});
  }

  deleteContact(email: String): Observable<any> {
    let headers = new Headers({ 'Content-Type': 'application/json' });
    let options = new RequestOptions({ headers: headers });
    return this.http.delete('http://localhost:8080/addressbook/contact/'+email, options).map(data => {return <ApiError> data.json()});
  }

}
