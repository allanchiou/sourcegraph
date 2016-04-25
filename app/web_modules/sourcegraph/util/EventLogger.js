// @flow weak

import React from "react";
import Dispatcher from "sourcegraph/Dispatcher";
import deepFreeze from "sourcegraph/util/deepFreeze";
import context from "sourcegraph/app/context";
import type {SiteConfig} from "sourcegraph/app/siteConfig";
import {getViewName, getRoutePattern} from "sourcegraph/app/routePatterns";
import type {Route} from "react-router";
import * as DashboardActions from "sourcegraph/dashboard/DashboardActions";
import * as UserActions from "sourcegraph/user/UserActions";
import * as DefActions from "sourcegraph/def/DefActions";
import UserStore from "sourcegraph/user/UserStore";

export const EventLocation = {
	Login: "Login",
	Signup: "Signup",
	Dashboard: "Dashboard",
	DefPopup: "DefPopup",
};

export class EventLogger {
	_amplitude: any;
	_intercomSettings: any;
	events: Array<any>;
	userProperties: Array<any>;
	intercomProperties: Array<any>;
	intercomEvents: Array<any>;
	isUserAgentBot: bool;
	_dispatcherToken: any;
	_siteConfig: ?SiteConfig;

	constructor() {
		this._amplitude = null;
		this._intercomSettings = null;

		this.events = deepFreeze([]);
		this.userProperties = deepFreeze([]);
		this.intercomProperties = deepFreeze([]);
		this.intercomEvents = deepFreeze([]);

		// Listen for all Stores dispatches.
		// You must separately log "frontend" actions of interest,
		// with the relevant event properties.
		this._dispatcherToken = Dispatcher.Stores.register(this.__onDispatch.bind(this));
	}

	// reset() receives any event data which is buffered
	// during server-side rendering; this data will
	// be flushed after the first call to
	// init() in the browser.
	reset(data) {
		this.events = deepFreeze(data && data.events ? data.events : this.events);
		this.userProperties = deepFreeze(data && data.userProperties ? data.userProperties : this.userProperties);
		this.intercomProperties = deepFreeze(data && data.intercomProperties ? data.intercomProperties : this.intercomProperties);
		this.intercomEvents = deepFreeze(data && data.intercomEvents ? data.intercomEvents : this.intercomEvents);
	}
	toJSON() {
		return {
			events: this.events,
			userProperties: this.userProperties,
			intercomProperties: this.intercomProperties,
			intercomEvents: this.intercomEvents,
		};
	}

	setSiteConfig(siteConfig: SiteConfig) {
		this._siteConfig = siteConfig;
	}

	// Loads the Amplitude JavaScript SDK if this
	// code is run in the browser (i.e. not with node
	// when doing server-side rendering.) If any events
	// have been buffered, it will flush them immediately.
	// If you do not call init() or it is run on the server,
	// any subequent calls to logEvent or setUserProperty
	// will be buffered.
	init() {
		const user = UserStore.activeUser();
		const emails = user && user.UID ? UserStore.emails.get(user.UID) : null;
		const primaryEmail = emails && !emails.Error ? emails.filter(e => e.Primary).map(e => e.Email)[0] : null;
		const authInfo = UserStore.activeAuthInfo();

		if (global.window && !this._amplitude) {
			this._amplitude = require("amplitude-js");

			if (!this._siteConfig) {
				throw new Error("EventLogger requires SiteConfig to be previously set using EventLogger.setSiteConfig before EventLogger can be initialized.");
			}

			let apiKey = "608f75cce80d583063837b8f5b18be54";
			if (this._siteConfig.buildVars.Version === "dev") {
				apiKey = "2b4b1117d1faf3960c81899a4422a222";
			} else {
				switch (this._siteConfig.appURL) {
				case "https://sourcegraph.com":
					apiKey = "e3c885c30d2c0c8bf33b1497b17806ba";
					break;
				case "https://staging.sourcegraph.com":
				case "https://staging2.sourcegraph.com":
				case "https://staging3.sourcegraph.com":
				case "https://staging4.sourcegraph.com":
					apiKey = "903f9390c3eefd5651853cf8dbd9d363";
					break;
				default:
					break;
				}
			}

			this._amplitude.init(apiKey, user ? user.Login : null, {
				includeReferrer: true,
			});

			if (user && user.RegisteredAt) {
				this.setUserProperty("registered_at", new Date(user.RegisteredAt).toDateString());
			}
			if (primaryEmail) {
				this.setUserProperty("email", primaryEmail);
			}
		}

		if (global.window) {
			this._intercomSettings = window.intercomSettings;
		}
		if (this._intercomSettings && user && authInfo) {
			this.setIntercomProperty("name", user.Name);
			if (primaryEmail) this.setIntercomProperty("email", primaryEmail);
			this.setIntercomProperty("user_id", user.UID.toString());
			this.setIntercomProperty("user_hash", authInfo.IntercomHash);
			this.setIntercomProperty("created_at", new Date(user.RegisteredAt).getTime() / 1000);
			// $FlowHack
			Intercom("boot", this._intercomSettings); // eslint-disable-line no-undef
		}

		// FullStory
		if (global.FS && user) {
			const id = user.Email || user.Login;
			// $FlowHack
			FS.identify(id, { // eslint-disable-line no-undef
				displayName: user.Name,
				email: primaryEmail,
			});
		}

		this.isUserAgentBot = Boolean(context.userAgentIsBot);
		this._flush();
	}

	// Only flush events on the client, after a call to init().
	// Filter out bot / test user agents.
	_shouldFlushAmplitude() {
		return Boolean(this._amplitude) && !this.isUserAgentBot;
	}
	_shouldFlushIntercom() {
		return Boolean(this._intercomSettings) && !this.isUserAgentBot;
	}

	// sets current user's properties
	setUserProperty(property, value) {
		if (!this._shouldFlushAmplitude()) {
			this.userProperties = deepFreeze(this.userProperties.concat([[property, value]]));
		} else {
			this._amplitude.identify(new this._amplitude.Identify().set(property, value));
		}
	}

	// records events for the current user
	logEvent(eventName, eventProperties) {
		if (!this._shouldFlushAmplitude()) {
			this.events = deepFreeze(this.events.concat([[eventName, eventProperties]]));
		} else {
			this._amplitude.logEvent(eventName, eventProperties);
		}
	}

	logEventForPage(eventName, pageName, eventProperties) {
		if (!pageName) throw new Error("PageName must be defined");

		let props = eventProperties ? eventProperties : {};
		props["page_name"] = pageName;
		this.logEvent(eventName, props);
	}

	// sets current user's property value
	setIntercomProperty(property, value) {
		if (!this._shouldFlushIntercom()) {
			this.intercomProperties = deepFreeze(this.intercomProperties.concat([[property, value]]));
		} else {
			this._intercomSettings[property] = value;
		}
	}

	// records intercom events for the current user
	logIntercomEvent(eventName, eventProperties) {
		if (!this._shouldFlushIntercom()) {
			this.intercomEvents = deepFreeze(this.intercomEvents.concat([[eventName, eventProperties]]));
		} else {
			window.Intercom("trackEvent", eventName, eventProperties);
		}
	}

	_flush() {
		if (this._shouldFlushAmplitude()) { // sanity check
			if (this.events) {
				for (let tuple of this.events) {
					this.logEvent(tuple[0], tuple[1]);
				}
				this.events = deepFreeze([]);
			}
			if (this.userProperties) {
				for (let tuple of this.userProperties) {
					this.setUserProperty(tuple[0], tuple[1]);
				}
				this.userProperties = deepFreeze([]);
			}
		}
		if (this._shouldFlushIntercom()) {
			if (this.intercomEvents) {
				for (let tuple of this.intercomEvents) {
					this.logIntercomEvent(tuple[0], tuple[1]);
				}
				this.intercomEvents = deepFreeze([]);
			}
			if (this.intercomProperties) {
				for (let tuple of this.intercomProperties) {
					this.setIntercomProperty(tuple[0], tuple[1]);
				}
				this.intercomProperties = deepFreeze([]);
			}
		}
	}

	__onDispatch(action) {
		switch (action.constructor) {
		case DashboardActions.RemoteReposFetched:
			if (action.data.RemoteRepos) {
				let orgs = {};
				for (let repo of action.data.RemoteRepos) {
					if (repo.OwnerIsOrg) orgs[repo.Owner] = true;
				}
				this.setUserProperty("orgs", Object.keys(orgs));
				this.setUserProperty("num_github_repos", action.data.RemoteRepos.length);
				this.setIntercomProperty("companies", Object.keys(orgs).map(org => ({id: `github_${org}`, name: org})));
			}
			break;

		case UserActions.SignupCompleted:
		case UserActions.LoginCompleted:
		case UserActions.LogoutCompleted:
		case UserActions.ForgotPasswordCompleted:
		case UserActions.ResetPasswordCompleted:
			if (action.email) {
				this.setUserProperty("email", action.email);
			}
			if (action.eventName) {
				this.logEvent(action.eventName, {error: Boolean(action.resp.Error)});
			}
			break;

		case DefActions.DefsFetched:
			if (action.eventName) {
				let eventProps = {
					query: action.query,
					overlay: action.overlay,
				};
				this.logEvent(action.eventName, eventProps);
			}
			break;
		default:
			// All dispatched actions to stores will automatically be tracked by the eventName
			// of the action (if set). Override this behavior by including another case above.
			if (action.eventName) {
				this.logEvent(action.eventName);
			}
			break;
		}

		this._flush(); // No need to __emitChange(); components need not be re-rendered.
	}
}

export default new EventLogger();

// withEventLoggerContext makes eventLogger accessible as this.context.eventLogger
// in the component's context.
export function withEventLoggerContext(eventLogger: EventLogger, Component: ReactClass): ReactClass {
	class WithEventLogger extends React.Component {
		static childContextTypes = {
			eventLogger: React.PropTypes.object,
		};

		constructor(props) {
			super(props);
			eventLogger.init();
		}

		getChildContext(): {eventLogger: EventLogger} {
			return {eventLogger};
		}

		render() {
			return <Component {...this.props} />;
		}
	}
	return WithEventLogger;
}

// withViewEventsLogged calls this.context.eventLogger.logEvent when the
// location's pathname changes.
export function withViewEventsLogged(Component: ReactClass): ReactClass {
	class WithViewEventsLogged extends React.Component { // eslint-disable-line react/no-multi-comp
		static propTypes = {
			routes: React.PropTypes.arrayOf(React.PropTypes.object),
			location: React.PropTypes.object,
		};

		static contextTypes = {
			router: React.PropTypes.object.isRequired,
			eventLogger: React.PropTypes.object.isRequired,
		};

		componentDidMount() {
			this._logView(this.props.routes, this.props.location);
		}

		componentWillReceiveProps(nextProps) {
			// Greedily log page views. Technically changing the pathname
			// may match the same "view" (e.g. interacting with the directory
			// tree navigations will change your URL,  but not feel like separate
			// page events). We will log any change in pathname as a separate event.
			// NOTE: this will not log separate page views when query string / hash
			// values are updated.
			if (this.props.location.pathname !== nextProps.location.pathname) {
				this._logView(nextProps.routes, nextProps.location);
			}
		}

		_logView(routes: Array<Route>, location: Location) {
			let eventProps = {
				referred_by_chrome_ext: false,
				url: location.pathname,
			};
			if (location.query && location.query["utm_source"] === "chromeext") {
				eventProps.referred_by_chrome_ext = true;
			}

			const viewName = getViewName(routes);
			if (viewName) {
				this.context.eventLogger.logEvent(viewName, eventProps);
			} else {
				this.context.eventLogger.logEvent("UnmatchedRoute", {
					...eventProps,
					pattern: getRoutePattern(routes),
				});
			}
		}

		render() { return <Component {...this.props} />; }
	}
	return WithViewEventsLogged;
}
