import { Component, OnInit } from '@angular/core';
import { AddressbookApiClientService } from '../addressbook-api-client.service';
import { Contact, ApiError } from '../shared/contact';

@Component({
  selector: 'app-addressbook-list',
  templateUrl: './addressbook-list.component.html',
  styleUrls: ['./addressbook-list.component.css']
})
export class AddressbookListComponent implements OnInit {
  contacts: Contact[]
  showTable = true;
  selectedContact: Contact
  selectedContactIndex: number
  isEdit: boolean
  constructor(private addressbookService: AddressbookApiClientService) { 
  }

  ngOnInit() {
    this.contacts = [];
    this.addressbookService.getContacts().subscribe(value =>this.contacts = value );
  } 
  
  deleteContact(email: String, index: number) {
    this.addressbookService.deleteContact(email).subscribe(apiError => {
      if(apiError.ErrorCode == "0") {
        this.contacts.splice(index, 1);
      }
    });
  }

  onAdded(contact: Contact) {
    console.log("called onAdded function");
    if (this.contacts == null) {
      this.contacts = [contact];
    } else {
      if (!this.isEdit) {
        this.contacts.push(contact);
      } else {
        console.log("editing contact")
        this.contacts[this.selectedContactIndex] = contact;
        this.isEdit = false;
      }
    }
    this.showTable = true;
  }

  setSelectedContact(contact: Contact, selectedContactIndex: number) {
    this.selectedContact = contact;
    this.selectedContactIndex = selectedContactIndex
    console.log(contact, selectedContactIndex);
    this.showTable = false;
    this.isEdit = true;
  }
}
