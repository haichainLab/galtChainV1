import { Injectable, NgZone } from '@angular/core';
import { Http } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { BehaviorSubject } from 'rxjs/BehaviorSubject';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class PriceService {
  readonly CMC_TICKER_ID = 1619;

  price: Subject<number> = new BehaviorSubject<number>(null);

  constructor(
    private http: Http,
    private ngZone: NgZone,
  ) {
    this.ngZone.runOutsideAngular(() => {
      Observable.timer(0, 10 * 60 * 1000).subscribe(() => {
        this.http.get(`http://samos.io/api/coin-price?token=galtcoin`)
          .map(response => response.json())
          .subscribe(response => this.ngZone.run(() => {
            this.price.next(response.data.quotes.CNY.price);
          }));
      });
    });
  }
}
