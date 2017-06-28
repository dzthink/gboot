/* Copyright 2017 dzthink
license that can be found in the LICENSE file.
*/

//define the interface that application should implement
package boot


type Application interface {
	Name() string   //name of the application
	Start(conf ConfigInterface) (err error) //application start hook
	Reload(conf ConfigInterface) (err error) //application reload hook
	Stop() (err error) //application stop hook
}
