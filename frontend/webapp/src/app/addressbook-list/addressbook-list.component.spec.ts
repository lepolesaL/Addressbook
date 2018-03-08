import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { AddressbookListComponent } from './addressbook-list.component';
import { PhonebookformComponent} from '../phonebookform/phonebookform.component';
import { ReactiveFormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';
import { AddressbookApiClientService } from '../addressbook-api-client.service';
import { Observable } from 'rxjs/observable';
import {of} from 'rxjs/observable/of';
import { CONTACTS } from '../shared/mockdata';
import { NO_ERRORS_SCHEMA } from '@angular/core';

describe('AddressbookListComponent', () => {
  let component: AddressbookListComponent;
  let fixture: ComponentFixture<AddressbookListComponent>;
  let apiService: jasmine.SpyObj<AddressbookApiClientService>;

  beforeEach(async(() => {
    const spy = jasmine.createSpyObj('AddressbookApiClientService', ['deleteContact']);
    TestBed.configureTestingModule({
      imports: [HttpModule, ReactiveFormsModule],
      declarations: [ AddressbookListComponent, PhonebookformComponent ],
      providers: [
        AddressbookApiClientService
      ]
    }).compileComponents();
    apiService = TestBed.get(AddressbookApiClientService);
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(AddressbookListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
     expect(component).toBeTruthy();
   });
});
