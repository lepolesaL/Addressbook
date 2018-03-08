import { Component, OnInit, EventEmitter, Input, Output  } from '@angular/core';
import { AddressbookApiClientService } from '../addressbook-api-client.service';

import {  FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Contact, ApiError } from '../shared/contact';

@Component({
  selector: 'app-phonebookform',
  templateUrl: './phonebookform.component.html',
  styleUrls: ['./phonebookform.component.css']
})
export class PhonebookformComponent implements OnInit {

  addressbookform: FormGroup;
  apiError:boolean
  apiErrorMessage:String
  @Input() contact: Contact
  @Input() isPut: boolean // is true when form is doing a PUT and false for post
  @Output() onAdded = new EventEmitter<Contact>();
  constructor(private fb: FormBuilder, private addressbookService: AddressbookApiClientService) { 
    this.createForm();
  }
  ngOnInit() {
    this.createForm();
  }

  createForm() {
    this.addressbookform = this.fb.group({
      name : ['', Validators.compose([Validators.required,  Validators.minLength(5),  Validators.maxLength(25)])],
      email: ['', Validators.compose([Validators.required, Validators.minLength(5), Validators.email])],
      phonenumber: ['', Validators.compose([Validators.required, Validators.minLength(10), Validators.maxLength(12)])],
      address: this.fb.group ({
        street: ['', Validators.compose([Validators.required, Validators.minLength(4), Validators.maxLength(50)])],
        city: ['', Validators.compose([Validators.required, Validators.minLength(4), Validators.maxLength(25)])],
        country: ['', Validators.compose([Validators.required, Validators.minLength(2), Validators.maxLength(2)])] //Change to country code
      }),
    });
  }

  ngOnChanges() {
    //console.log("value changed: ", this.contact, this.isPut);
    if(this.contact != null) {

      this.addressbookform.setValue({
        name : this.contact.name,
        email: this.contact.email,
        phonenumber: this.contact.phonenumber,
        address: this.contact.address
      });
    }
  }

  addContact({ value, valid }: { value: Contact, valid: boolean}) {
    console.log(value, valid);
    if(valid) {
      if(!this.isPut) {
        this.addressbookService.addContact(value).subscribe(apiError => {
          // check for success and then navigate to parent view
          if (apiError.ErrorCode == "0") {
            console.log("Sucessful added new contact");
            this.onAdded.emit(value);
          } else {
            this.apiError = true;
            this.apiErrorMessage = apiError.Message;
            console.log(this.apiErrorMessage);

          }
        });
      } else {
        this.addressbookService.editContact(value).subscribe(apiError => {
          // check for success and then navigate to parent view
          if (apiError.ErrorCode == "0") {
            console.log("Successful updating contact");
            this.onAdded.emit(value);
          } else {
              this.apiError = true;
              this.apiErrorMessage = apiError.Message;
          }
        });
      } 
    }
  }

  get name() { return this.addressbookform.get('name'); }

  get email() { return this.addressbookform.get('email'); }

  get phonenumber() { return this.addressbookform.get('phonenumber'); }

  get street() { return this.addressbookform.get('address').get('street'); }

  get city() { return this.addressbookform.get('address').get('city'); }

  get country() { return this.addressbookform.get('address').get('country'); }

}

