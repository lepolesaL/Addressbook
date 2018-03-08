import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PhonebookformComponent } from './phonebookform.component';
import { AddressbookApiClientService } from '../addressbook-api-client.service';
import { Contact } from '../shared/contact';
import { HttpModule, Http } from '@angular/http';
import { AppModule } from '../app.module';
import { ReactiveFormsModule } from '@angular/forms';
import { Observable } from 'rxjs/observable';
import {of} from 'rxjs/observable/of';


describe('PhonebookformComponent', () => {
  let component: PhonebookformComponent;
  let fixture: ComponentFixture<PhonebookformComponent>;
  let apiService: jasmine.SpyObj<AddressbookApiClientService>;

  beforeEach(async(() => {
    const spy = jasmine.createSpyObj('AddressbookApiClientService', ['addContact', 'editContact']);
    TestBed.configureTestingModule({
      imports: [HttpModule, ReactiveFormsModule],
      declarations: [ PhonebookformComponent ],
      providers: [
        { provide: AddressbookApiClientService, useValue: spy }
      ]
    })
    .compileComponents();
    apiService = TestBed.get(AddressbookApiClientService);
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PhonebookformComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should add Contact', () => {
    let value: Contact = {name: "name", email:"email", phonenumber:"012345", 
    address:{street:"stree", city:"city", country:"ZA"}};
    let valid: boolean = true;
    apiService.addContact.and.returnValue(of({errorCode:"0", Message: "Success"}));
    expect(component.addContact({value, valid})).toBeTruthy;
    expect(apiService.addContact.calls.count()).toBe(1, "Api client service called");
  });

  it('should not add contact', () => {
    let value: Contact = {name: "name", email:"email", phonenumber:"012345", 
    address:{street:"street", city:"city", country:"ZA"}};
    let valid: boolean = false;
    apiService.addContact.and.returnValue(<any>{});
    expect(component.addContact({value, valid})).toBeTruthy;
    expect(apiService.addContact.calls.count()).toBe(0, "Api client service is not called");
  });

  it('should edit contact', () => {
    component.isPut = true;
    let value: Contact = {name: "name", email:"email", phonenumber:"012345", 
    address:{street:"street", city:"city", country:"ZA"}};
    let valid: boolean = false;
    apiService.addContact.and.returnValue(<any>{});
    apiService.editContact.and.returnValue(of({errorCode:"0", Message:"Success"}));
    expect(component.addContact({value, valid})).toBeTruthy;
    expect(apiService.addContact.calls.count()).toBe(0, "Api client service is not called");
    expect(apiService.editContact.calls.count()).toBe(0, "Api client service is not called");
  });
});
