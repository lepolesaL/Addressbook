import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { AddressbookApiClientService } from './addressbook-api-client.service';
import { HttpModule } from '@angular/http';


import { AppComponent } from './app.component';
import { PhonebookformComponent } from './phonebookform/phonebookform.component';
import { ReactiveFormsModule } from '@angular/forms';
import { AddressbookListComponent } from './addressbook-list/addressbook-list.component';


@NgModule({
  declarations: [
    AppComponent,
    PhonebookformComponent,
    AddressbookListComponent
  ],
  imports: [
    BrowserModule,
    ReactiveFormsModule,
    HttpModule
  ],
  providers: [
    AddressbookApiClientService,
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
