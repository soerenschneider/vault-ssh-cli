# Changelog

## [1.9.0](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.8.3...v1.9.0) (2025-01-01)


### Features

* support Vault Kubernetes auth method and writing/reading from Kubernetes secrets ([3603fe9](https://github.com/soerenschneider/vault-ssh-cli/commit/3603fe9d2eeac9d90a47eda6105fdf4d040b5b7c))


### Bug Fixes

* **deps:** bump github.com/hashicorp/vault/api from 1.14.0 to 1.15.0 ([#190](https://github.com/soerenschneider/vault-ssh-cli/issues/190)) ([945532f](https://github.com/soerenschneider/vault-ssh-cli/commit/945532f84b5a44c40a78866f5727cc7ad1b9d1e9))
* **deps:** bump github.com/prometheus/client_golang ([#199](https://github.com/soerenschneider/vault-ssh-cli/issues/199)) ([0a36847](https://github.com/soerenschneider/vault-ssh-cli/commit/0a368477d9ef51de7f99c2cf51f93b93fba2e8af))
* **deps:** bump golang.org/x/crypto from 0.26.0 to 0.31.0 ([#208](https://github.com/soerenschneider/vault-ssh-cli/issues/208)) ([60f1102](https://github.com/soerenschneider/vault-ssh-cli/commit/60f11023d5417400ace78063a1fd84730539aba2))
* **deps:** bump golang.org/x/term from 0.24.0 to 0.27.0 ([#206](https://github.com/soerenschneider/vault-ssh-cli/issues/206)) ([8dfa5b8](https://github.com/soerenschneider/vault-ssh-cli/commit/8dfa5b87de2d3523ee0b74da4adf0c2f35ca063d))

## [1.8.3](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.8.2...v1.8.3) (2024-09-09)


### Bug Fixes

* **deps:** bump github.com/go-playground/validator/v10 ([#185](https://github.com/soerenschneider/vault-ssh-cli/issues/185)) ([da821bb](https://github.com/soerenschneider/vault-ssh-cli/commit/da821bbf98648ce3e63bdf5e3f76793d9c681e4c))
* **deps:** bump github.com/prometheus/client_golang ([#186](https://github.com/soerenschneider/vault-ssh-cli/issues/186)) ([056de3a](https://github.com/soerenschneider/vault-ssh-cli/commit/056de3a5a9519db44b42e2f8f1021f3711bc2077))
* **deps:** bump github.com/prometheus/common from 0.54.0 to 0.59.1 ([#187](https://github.com/soerenschneider/vault-ssh-cli/issues/187)) ([39a8278](https://github.com/soerenschneider/vault-ssh-cli/commit/39a82784e843d7e2a2e204014126e4f231d63a01))
* **deps:** bump github.com/spf13/cobra from 1.8.0 to 1.8.1 ([#182](https://github.com/soerenschneider/vault-ssh-cli/issues/182)) ([9c4a47d](https://github.com/soerenschneider/vault-ssh-cli/commit/9c4a47d593dcad96d926eab475a6ccc391954cee))
* **deps:** bump golang.org/x/term from 0.23.0 to 0.24.0 ([#184](https://github.com/soerenschneider/vault-ssh-cli/issues/184)) ([0825665](https://github.com/soerenschneider/vault-ssh-cli/commit/08256653b31dc0d956213f8bf634a73039e8f4c8))
* fix vault auth ([be3aa87](https://github.com/soerenschneider/vault-ssh-cli/commit/be3aa8750b2ea5a1d8e2ec28673bdcd14101d0a1))

## [1.8.2](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.8.1...v1.8.2) (2024-08-24)


### Bug Fixes

* **deps:** bump github.com/cenkalti/backoff/v3 from 3.0.0 to 3.2.2 ([#148](https://github.com/soerenschneider/vault-ssh-cli/issues/148)) ([c2c3270](https://github.com/soerenschneider/vault-ssh-cli/commit/c2c3270189c34c36e0c418407c0ed0937cf1dd4a))
* **deps:** bump github.com/hashicorp/go-retryablehttp ([#159](https://github.com/soerenschneider/vault-ssh-cli/issues/159)) ([559a7b0](https://github.com/soerenschneider/vault-ssh-cli/commit/559a7b033de69581d7e820efd1fc143149fef574))
* **deps:** bump golang.org/x/crypto from 0.23.0 to 0.26.0 ([#172](https://github.com/soerenschneider/vault-ssh-cli/issues/172)) ([8fe6254](https://github.com/soerenschneider/vault-ssh-cli/commit/8fe625416b474a9dcf7f384bc7759373c9fcd167))
* **deps:** bump golang.org/x/sys from 0.20.0 to 0.24.0 ([#173](https://github.com/soerenschneider/vault-ssh-cli/issues/173)) ([199abf9](https://github.com/soerenschneider/vault-ssh-cli/commit/199abf959fa4d5c497cb53fcf0a504a8ccfdfc47))
* **deps:** bump golang.org/x/term from 0.20.0 to 0.23.0 ([#174](https://github.com/soerenschneider/vault-ssh-cli/issues/174)) ([95e178e](https://github.com/soerenschneider/vault-ssh-cli/commit/95e178e09f5706b0c6993e771d1bfa813bf5b801))
* fix duplicate method ([9159496](https://github.com/soerenschneider/vault-ssh-cli/commit/9159496326127d6bce5282e076be43412a8a5380))
* remove duplicate login call ([788c724](https://github.com/soerenschneider/vault-ssh-cli/commit/788c7248c79e8c9f74c69d619956b94bd59fa2ee))

## [1.8.1](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.8.0...v1.8.1) (2024-06-10)


### Bug Fixes

* expand config file ([299296f](https://github.com/soerenschneider/vault-ssh-cli/commit/299296fbcad88e94335e1c3eedc960398f5233b4))
* set viper default for retries ([534b213](https://github.com/soerenschneider/vault-ssh-cli/commit/534b2138a9af05a9d9bbb25ff35408120642c3fe))

## [1.8.0](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.7.2...v1.8.0) (2024-06-03)


### Features

* increase resilience by using exponential backoff ([7f7588b](https://github.com/soerenschneider/vault-ssh-cli/commit/7f7588b6d504b4974c514c9371f40659078da7f7))
* make amount of retries configurable ([a8f7139](https://github.com/soerenschneider/vault-ssh-cli/commit/a8f7139ab930675fb5f91d4ae5195ae5b5368ce5))


### Bug Fixes

* **deps:** bump github.com/go-playground/validator/v10 ([78efacb](https://github.com/soerenschneider/vault-ssh-cli/commit/78efacbe49739e2b56a00269d9101b658d24938b))
* **deps:** bump github.com/prometheus/client_golang ([65b0637](https://github.com/soerenschneider/vault-ssh-cli/commit/65b0637200befa6a395c183c732585e7e8911bc5))
* **deps:** bump github.com/prometheus/common from 0.53.0 to 0.54.0 ([c2c79d7](https://github.com/soerenschneider/vault-ssh-cli/commit/c2c79d70e13e258d2c566d0548e6956bcb78ed19))
* **deps:** bump github.com/rs/zerolog from 1.31.0 to 1.33.0 ([0d3afb6](https://github.com/soerenschneider/vault-ssh-cli/commit/0d3afb63276733982a04c49d3ba0c6797be4c936))
* **deps:** bump github.com/spf13/viper from 1.17.0 to 1.19.0 ([b59e455](https://github.com/soerenschneider/vault-ssh-cli/commit/b59e4557e1e153af661c08951178a03070b335de))

## [1.7.2](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.7.1...v1.7.2) (2024-06-02)


### Bug Fixes

* don't return error when cert file does not exist, yet ([7634a9a](https://github.com/soerenschneider/vault-ssh-cli/commit/7634a9a26da155217a778741b60e4ffd48801dc1))
* omit validation of field when empty ([7bc2a85](https://github.com/soerenschneider/vault-ssh-cli/commit/7bc2a85001b0b05942189bbde67954478e80b3f6))
* use the same pod impl for every all commands ([fa58a9b](https://github.com/soerenschneider/vault-ssh-cli/commit/fa58a9b74aa849217ea618da709e0fec116f1750))

## [1.7.1](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.7.0...v1.7.1) (2024-05-29)


### Bug Fixes

* **deps:** bump github.com/go-jose/go-jose/v3 from 3.0.0 to 3.0.3 ([601fcc5](https://github.com/soerenschneider/vault-ssh-cli/commit/601fcc59958441e54946cd3b8474c82d0c7a5f9d))
* **deps:** bump github.com/go-playground/validator/v10 ([83ca68f](https://github.com/soerenschneider/vault-ssh-cli/commit/83ca68fbe8a2261ca71a4443dd34d2f6dd58ae1d))
* **deps:** bump github.com/prometheus/common from 0.44.0 to 0.53.0 ([9f2f585](https://github.com/soerenschneider/vault-ssh-cli/commit/9f2f585e44d5b93e31630b231b19020d079e6b67))
* **deps:** bump github.com/spf13/cobra from 1.7.0 to 1.8.0 ([f64952b](https://github.com/soerenschneider/vault-ssh-cli/commit/f64952be5cad0cd2edf6cec5ca1fa95da0d75e10))
* **deps:** bump golang.org/x/crypto from 0.21.0 to 0.23.0 ([bc838da](https://github.com/soerenschneider/vault-ssh-cli/commit/bc838dac168c890215877a61e5dd8b3d3e189ec0))
* **deps:** bump golang.org/x/net from 0.19.0 to 0.23.0 ([bf3d411](https://github.com/soerenschneider/vault-ssh-cli/commit/bf3d411c19d02826c10e4795800b7278559c9766))
* **deps:** bump golang.org/x/sys from 0.18.0 to 0.20.0 ([c0a5df1](https://github.com/soerenschneider/vault-ssh-cli/commit/c0a5df19fdf26261fd64415ca2b3c82aad076eae))
* **deps:** bump google.golang.org/protobuf from 1.31.0 to 1.33.0 ([7b285e1](https://github.com/soerenschneider/vault-ssh-cli/commit/7b285e1737a15f7b1e01c8c2ecb4465ebdccb469))
* fix dumping metrics with new prometheus version ([47f9b6c](https://github.com/soerenschneider/vault-ssh-cli/commit/47f9b6c8698c9d6ce2c5530cbea56c8dc9f06a12))
* use correct method ([6cdc1a7](https://github.com/soerenschneider/vault-ssh-cli/commit/6cdc1a7613d78fcd6721cdc7886da953fd4629d9))

## [1.7.0](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.6.2...v1.7.0) (2024-05-29)


### Features

* add support for setting ttl and principals ([2e7b6d1](https://github.com/soerenschneider/vault-ssh-cli/commit/2e7b6d1d61560cda540ecc97b6c30e205e0f16b7))


### Bug Fixes

* 121: Provide SshMountPath instead of VaultRole for readCACert: ([b58b26c](https://github.com/soerenschneider/vault-ssh-cli/commit/b58b26c38134853a9ab69cf2065b6c3c53b674b0))
* expand all files ([3781095](https://github.com/soerenschneider/vault-ssh-cli/commit/3781095158d1d39144722a3dd8a7531349f0b829))

## [1.6.2](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.6.1...v1.6.2) (2023-12-24)


### Bug Fixes

* fix print version cmd ([6ff4776](https://github.com/soerenschneider/vault-ssh-cli/commit/6ff4776552e437c0ca7341a8aa81d8257d4c7c0e))

## [1.6.1](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.6.0...v1.6.1) (2023-12-23)


### Bug Fixes

* **deps:** bump github.com/prometheus/client_golang ([#88](https://github.com/soerenschneider/vault-ssh-cli/issues/88)) ([ca303d8](https://github.com/soerenschneider/vault-ssh-cli/commit/ca303d8ffe7692dc975677892aec64c86a9bc64d))
* **deps:** bump github.com/spf13/viper from 1.16.0 to 1.17.0 ([#91](https://github.com/soerenschneider/vault-ssh-cli/issues/91)) ([8c5a196](https://github.com/soerenschneider/vault-ssh-cli/commit/8c5a196d5d88b894119c23dfbb134f02a905992e))
* **deps:** bump golang.org/x/crypto from 0.13.0 to 0.14.0 ([#90](https://github.com/soerenschneider/vault-ssh-cli/issues/90)) ([f826ed4](https://github.com/soerenschneider/vault-ssh-cli/commit/f826ed48f2cd88ad6f903dff747058ddea25c04d))
* **deps:** bump golang.org/x/net from 0.15.0 to 0.17.0 ([#93](https://github.com/soerenschneider/vault-ssh-cli/issues/93)) ([44807d7](https://github.com/soerenschneider/vault-ssh-cli/commit/44807d74dfd4feb61c50b64d3e71a182d5562e2a))
* only dump metrics defined by the tool itself ([a32ecfe](https://github.com/soerenschneider/vault-ssh-cli/commit/a32ecfe23cc58454e8fe2a71f879c9a05076543d))

## [1.6.0](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.5.1...v1.6.0) (2023-09-22)


### Features

* set global loglevel, auto detect whether running on terminal ([beba023](https://github.com/soerenschneider/vault-ssh-cli/commit/beba023a057858abd0be7459668cbb79ddcbc95d))

## [1.5.1](https://github.com/soerenschneider/vault-ssh-cli/compare/v1.5.0...v1.5.1) (2023-08-07)


### Bug Fixes

* expand user paths only at single point ([fe2b311](https://github.com/soerenschneider/vault-ssh-cli/commit/fe2b3114e9b1d4db12ca700b463ded1f16d9022a))

## [1.5.0](https://github.com/soerenschneider/ssh-key-signer/compare/v1.4.1...v1.5.0) (2023-08-02)


### Miscellaneous Chores

* release 1.5.0 ([64ada5a](https://github.com/soerenschneider/ssh-key-signer/commit/64ada5a49845717b0aeb9cd0ddbd1dee30e56207))

## [1.4.1](https://github.com/soerenschneider/ssh-key-signer/compare/v1.4.0...v1.4.1) (2023-05-02)


### Bug Fixes

* respect '--force-new-signature' flag ([8e3a62b](https://github.com/soerenschneider/ssh-key-signer/commit/8e3a62bcd119cff80738fea348b37bbe15f8f0cb))

## [1.4.0](https://github.com/soerenschneider/ssh-key-signer/compare/v1.3.1...v1.4.0) (2023-02-09)


### Features

* sign user certs ([#38](https://github.com/soerenschneider/ssh-key-signer/issues/38)) ([2f06648](https://github.com/soerenschneider/ssh-key-signer/commit/2f066481eb34005a53cc5d5025529f4ad7149572))

## [1.3.1](https://github.com/soerenschneider/ssh-key-signer/compare/v1.3.0...v1.3.1) (2023-01-24)


### Miscellaneous Chores

* release 1.3.1 ([ee4fe98](https://github.com/soerenschneider/ssh-key-signer/commit/ee4fe98c7eca86543c68097a96e03db78861211b))
* release 1.3.3 ([6c234cb](https://github.com/soerenschneider/ssh-key-signer/commit/6c234cb41f0a0ff98a1b44033a9785612473b6d3))

## [1.3.0](https://github.com/soerenschneider/ssh-key-signer/compare/v1.2.1...v1.3.0) (2023-01-24)


### Features

* Add subcmd for reading ca ([#33](https://github.com/soerenschneider/ssh-key-signer/issues/33)) ([8a1caa2](https://github.com/soerenschneider/ssh-key-signer/commit/8a1caa2fa1def937b21ce7fba982733ab3d6218b))

## [1.2.1](https://github.com/soerenschneider/ssh-key-signer/compare/v1.2.0...v1.2.1) (2023-01-06)


### Bug Fixes

* do not write 'go_' prefixed metrics to file ([d270167](https://github.com/soerenschneider/ssh-key-signer/commit/d27016788bbb1028f2f86aaa5273179846efd5bc))

## [1.2.0](https://www.github.com/soerenschneider/ssh-key-signer/compare/v1.1.0...v1.2.0) (2022-02-08)


### Features

* add version info, log on startup ([ce23d40](https://www.github.com/soerenschneider/ssh-key-signer/commit/ce23d40fb5ed7de5a2637718ff16030f69aab4c7))
* split functionality into commands ([1dbdbde](https://www.github.com/soerenschneider/ssh-key-signer/commit/1dbdbde75b1fa88b44ad9f61d0dc93c7e98433e5))

## [1.1.0](https://www.github.com/soerenschneider/ssh-key-signer/compare/v1.0.0...v1.1.0) (2022-02-03)


### Features

* expose new metrics ([79eb199](https://www.github.com/soerenschneider/ssh-key-signer/commit/79eb19948bd3688a61e2f0f0a01e7993c98954bc))

## 1.0.0 (2022-02-02)


### Bug Fixes

* Fix incorrect check for AppRole auth data ([0320284](https://www.github.com/soerenschneider/ssh-key-signer/commit/032028422be048f0c874e23c969db802ea929dd2))
* Include 'cert_type' parameter ([5751678](https://www.github.com/soerenschneider/ssh-key-signer/commit/5751678afef2a69b62048ab4c76ee535748229ec))
* Set AppRole mount path ([96b587f](https://www.github.com/soerenschneider/ssh-key-signer/commit/96b587f3607ee972fb9d150f3fa8cff3bd3ff937))
