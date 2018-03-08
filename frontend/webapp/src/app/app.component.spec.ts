import { TestBed, async } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { AddressbookListComponent } from './addressbook-list/addressbook-list.component';
import { PhonebookformComponent } from './phonebookform/phonebookform.component';
import { ReactiveFormsModule } from '@angular/forms';
import { AddressbookApiClientService } from './addressbook-api-client.service';
import { HttpModule } from '@angular/http';
describe('AppComponent', () => {
  let app: AppComponent
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports:[ReactiveFormsModule, HttpModule],
      declarations: [
        AppComponent, AddressbookListComponent
      ],
      providers: [AddressbookApiClientService]
    }).overrideComponent(AddressbookListComponent, {
      set: {
        selector: 'app-addressbook-list',
        template: '<span>address book list</span>'
      }
    })
    .compileComponents();
  }));
  beforeEach(() => {
    const fixture = TestBed.createComponent(AppComponent);
    app = fixture.debugElement.componentInstance;
    fixture.detectChanges();
  });
  it('should create the app', () => {
    expect(app).toBeTruthy();
  });
});
